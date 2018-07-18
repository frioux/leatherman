leatherman.xz: leatherman
	xz leatherman

leatherman: *.go
	go get github.com/golang/lint
	go get -t ./...
	PATH="$(PATH):$(shell go env GOPATH)/bin"
	golint -set_exit_status ./...
	go vet ./...
	go test
	go build -ldflags "-X main.version=$(TRAVIS_COMMIT)"
	strip leatherman
