leatherman: *.go
	go get -t ./...
	go test
	go build -ldflags "-X main.version=`git describe --dirty --all`"
	strip leatherman
	xz leatherman
