leatherman: *.go
	go get -t ./...
	go test
	go build -ldflags "-X main.version=`git log -1 HEAD --pretty=%H`"
	strip leatherman
	xz leatherman
