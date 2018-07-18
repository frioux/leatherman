leatherman.xz: leatherman
	xz leatherman

leatherman: *.go
	go get -v -t ./...
	golint -set_exit_status ./...
	go vet ./...
	go test
	go build -ldflags "-X main.version=$(TRAVIS_COMMIT)"
	strip leatherman
