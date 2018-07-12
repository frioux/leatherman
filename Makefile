leatherman.xz: leatherman
	xz leatherman

leatherman: *.go
	go get -t ./...
	go test
	go build -ldflags "-X main.version=$(TRAVIS_COMMIT)"
	strip leatherman
