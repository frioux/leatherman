leatherman: *.go
	go get -t ./...
	go test
	go build
	strip leatherman
