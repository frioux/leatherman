package main

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/elb"
	"net"
	"os"
	"strings"
	"sync"
)

type ipLookup struct {
	name   string
	result string
	err    error
}

var verbose bool

func allRegions(cfg aws.Config) ([]string, error) {
	svc := ec2.New(cfg)

	ret := []string{}

	resp, err := svc.DescribeRegionsRequest(nil).Send()
	if err != nil {
		return []string{"us-west-1", "us-east-1", "us-west-2"}, err
	}

	for _, region := range resp.Regions {
		ret = append(ret, *region.RegionName)
	}
	return ret, nil
}

func Ec2ResourceForIp(args []string) {
	flags := flag.NewFlagSet("ec2-resource-for-ip", flag.ExitOnError)

	flags.BoolVar(&verbose, "verbose", false, "Never stop talking")
	e := flags.Parse(args[1:])
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}

	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}
	cfg.Region = endpoints.UsWest2RegionID

	regions, err := allRegions(cfg)
	if verbose {
		fmt.Println("Checking", regions)
	}
	if err != nil {
		if verbose {
			fmt.Println(err)
		}
	}

	ips := flags.Args()

	out := make(chan ipLookup)
	errC := make(chan error)

	done := make(chan bool)

	var wg sync.WaitGroup
	var workerWg sync.WaitGroup

	wg.Add(len(ips))
	workerWg.Add(1)
	foundIps := make(map[string]bool)
	for _, ip := range ips {
		foundIps[ip] = false
	}

	go func() {
		workerWg.Wait()
		close(done)
	}()

	go func() {
		wg.Wait()
		close(done)
	}()

	go func() {
		for err := range errC {
			if err != nil {
				if verbose {
					fmt.Println(err)
				}
			}
		}
	}()

	go func() {
		for i := range out {
			if i.err != nil {
				if verbose {
					fmt.Println(i.err)
				}
			} else {
				delete(foundIps, i.name)
				fmt.Print(i.name + ":\n" + i.result)
			}
			// XXX: need to verify that the thing found actually deleted a
			// record, or something
			wg.Done()
		}
	}()

	for _, region := range regions {
		region := region

		workerWg.Add(1)
		go func() {
			errC <- ec2InstancePublic(region, cfg, ips, out)
			workerWg.Done()
		}()

		workerWg.Add(1)
		go func() {
			errC <- ec2InstancePrivate(region, cfg, ips, out)
			workerWg.Done()
		}()

		workerWg.Add(1)
		go func() {
			errC <- findEIP(region, cfg, ips, out)
			workerWg.Done()
		}()

		workerWg.Add(1)
		go func() {
			errC <- findELB(region, cfg, ips, errC, out)
			workerWg.Done()
		}()
	}
	workerWg.Done()

	for range done {

	}

	keys := make([]string, len(foundIps))
	i := 0
	for k := range foundIps {
		keys[i] = k
		i++
	}

	found, err := unknown(keys)
	if err != nil {
		if verbose {
			fmt.Println(err)
		}
	} else {
		for ip, str := range found {
			delete(foundIps, ip)
			fmt.Print(ip + ":\n" + str)
		}
	}

	for ip := range foundIps {
		fmt.Printf("%s:\n", ip)
	}
}

func findELB(region string, cfg aws.Config, ips []string, errC chan error, out chan ipLookup) error {
	cfg.Region = region
	svc := elb.New(cfg)

	// map of ip to elb-id
	lookup := make(map[string]string)

	resp, err := svc.DescribeLoadBalancersRequest(nil).Send()

	if err != nil {
		return err
	}

	for _, lb := range resp.LoadBalancerDescriptions {
		ips, err := net.LookupIP(*lb.DNSName)

		// This happens all the time; do not early exit
		if err != nil {
			errC <- err
			continue
		}

		name := *lb.LoadBalancerName
		for _, ip := range ips {
			lookup[ip.String()] = name
		}
	}

	for _, ip := range ips {
		if name, ok := lookup[ip]; ok {
			out <- ipLookup{
				name: ip,
				result: fmt.Sprintf(
					"  type: elb\n"+
						"  region: %s\n"+
						"  name: %s\n", region, name),
			}
		}
	}

	return nil
}

func findEIP(region string, cfg aws.Config, ips []string, out chan ipLookup) error {
	cfg.Region = region
	svc := ec2.New(cfg)

	params := &ec2.DescribeAddressesInput{
		Filters: []ec2.Filter{
			{
				Name:   aws.String("public-ip"),
				Values: ips,
			},
		},
	}

	resp, err := svc.DescribeAddressesRequest(params).Send()
	if err != nil {
		return err
	}

	for _, address := range resp.Addresses {
		id := address.AllocationId
		if id == nil {
			continue
		}
		out <- ipLookup{
			name: *address.PublicIp,
			result: fmt.Sprintf(
				"  type: eip\n"+
					"  region: %s\n"+
					"  id: %s\n", region, *id),
		}
	}
	return nil
}

func toptr(ip string) string {
	parts := strings.Split(ip, ".")

	ret := ""
	for i := len(parts) - 1; i > 0; i-- {
		ret += parts[i] + "."
	}
	ret += parts[0] + ".in-addr.arpa"

	return ret
}

func unknown(ips []string) (map[string]string, error) {
	ret := make(map[string]string)
	for _, ip := range ips {
		ptrs, err := net.LookupAddr(ip)
		if err != nil {
			if verbose {
				fmt.Println(err)
			}
			ptrs = []string{""}
		}
		ret[ip] = fmt.Sprintf(
			"  type: unknown\n"+
				"  ptr: %s\n", ptrs[0])
	}
	return ret, nil
}

func getEC2Name(i ec2.Instance) string {
	for _, tag := range i.Tags {
		if *tag.Key == "Name" {
			return *tag.Value
		}
	}
	return ""
}

func ec2InstancePublic(region string, cfg aws.Config, ips []string, out chan ipLookup) error {
	cfg.Region = region
	svc := ec2.New(cfg)

	params := &ec2.DescribeInstancesInput{
		Filters: []ec2.Filter{
			{
				Name:   aws.String("ip-address"),
				Values: ips,
			},
		},
	}
	resp, err := svc.DescribeInstancesRequest(params).Send()

	if err != nil {
		return err
	}

	for _, res := range resp.Reservations {
		for _, instance := range res.Instances {
			out <- ipLookup{
				name: *instance.PublicIpAddress,
				result: fmt.Sprintf(
					"  type: ec2_instance\n"+
						"  region: %s\n"+
						"  id: %s\n"+
						"  name: %s\n",
					region, *instance.InstanceId, getEC2Name(instance)),
			}
		}
	}
	return nil
}

func ec2InstancePrivate(region string, cfg aws.Config, ips []string, out chan ipLookup) error {
	cfg.Region = region
	svc := ec2.New(cfg)

	params := &ec2.DescribeInstancesInput{
		Filters: []ec2.Filter{
			{
				Name:   aws.String("private-ip-address"),
				Values: ips,
			},
		},
	}
	resp, err := svc.DescribeInstancesRequest(params).Send()

	if err != nil {
		return err
	}

	for _, res := range resp.Reservations {
		for _, instance := range res.Instances {
			out <- ipLookup{
				name: *instance.PrivateIpAddress,
				result: fmt.Sprintf(
					"  type: ec2_instance\n"+
						"  region: %s\n"+
						"  id: %s\n"+
						"  name: %s\n",
					region, *instance.InstanceId, getEC2Name(instance)),
			}
		}
	}
	return nil
}
