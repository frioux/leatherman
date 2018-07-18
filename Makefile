leatherman.xz: leatherman
	xz leatherman

leatherman: *.go
	go get -t ./...
	go vet ./...
	go test
	go build -ldflags "-X main.version=$(TRAVIS_COMMIT)"
	strip leatherman

lint:
	golint -set_exit_status ./...
