leatherman: *.go
	go get -t ./...
	go test
	go build
	bzip2 leatherman
