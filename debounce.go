package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"time"
)

func cat(c chan string, e chan error, quit chan struct{}) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		c <- scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		e <- err
	}
	quit <- struct{}{}
}

func Debounce(args []string) {
	var timeoutSeconds float64
	var begin, end, h, help bool

	flags := flag.NewFlagSet("debounce", flag.ExitOnError)

	flags.Float64Var(&timeoutSeconds, "lockoutTime", 1, "amount of time between output")
	flags.BoolVar(&begin, "leadingEdge", false, "trigger at leading edge of cycle")
	flags.BoolVar(&end, "trailingEdge", true, "trigger at trailing edge of cycle")
	flags.BoolVar(&h, "h", false, "help for debounce")
	flags.BoolVar(&help, "help", false, "help for debounce")

	e := flags.Parse(args[1:])
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}

	if h || help {
		fmt.Println("\n" +
			" debounce          [--leadingEdge] [--trailingEdge] [--lockoutTime 2]\n" +
			"                   [-h|--help]\n" +
			"\n" +
			"    --leadingEdge   pass this flag to output at the leading edge of a cycle\n" +
			"                    (off by default)\n" +
			"    --trailingEdge  pass this flag to output at the trailing edge of a cycle\n" +
			"                    (on by default, pass false to disable)\n" +
			"    --lockoutTime   set the lockout time in seconds, default is 1 second\n" +
			"\n" +
			"    -h --help       print usage message and exit\n" +
			"\n" +
			"\n" +
			"debounce creates cycles based on the lockout time.  The cycle\n" +
			"starts on the first line sent and stops after no lines are sent\n" +
			"within a period of the lockout time\n" +
			"\n" +
			"\n" +
			"The following would run tests after a second of 'silence' after a\n" +
			"save\n" +
			"\n" +
			" inotifywait -mr -e modify,move . | debounce | xargs -i{} make test\n" +
			"",
		)
		return
	}

	c := make(chan string)
	quit := make(chan struct{})
	err := make(chan error)

	go cat(c, err, quit)

	for {
		var x string
		shouldPrint := false
		select {
		case x = <-c:
			shouldPrint = true
			if begin {
				shouldPrint = false
				fmt.Println(x)
			}
		case x := <-err:
			fmt.Fprintln(os.Stderr, "reading standard input:", x)
		case <-quit:
			return
		}
		timeout := time.After(time.Duration(timeoutSeconds) * time.Second)
	InnerLoop:
		for {
			select {
			case x = <-c:
				shouldPrint = true
				timeout = time.After(time.Duration(timeoutSeconds) * time.Second)
			case x := <-err:
				fmt.Fprintln(os.Stderr, "reading standard input:", x)
			case <-quit:
				return
			case <-timeout:
				if end && shouldPrint {
					fmt.Println(x)
				}
				break InnerLoop
			}
		}
	}
}
