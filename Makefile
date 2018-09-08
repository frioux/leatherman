VERSION := $(shell git rev-parse HEAD || echo $TRAVIS_COMMIT)
WHEN := $(shell date)
WHO := $(shell whoami)
WHERE := $(shell hostname -f)

leatherman.xz: leatherman
	xz leatherman

leatherman: *.go
	go get -t ./...
	go get github.com/golang/lint/golint
	golint -set_exit_status ./...
	go vet ./...
	TZ=America/Los_Angeles go test
	go build -ldflags "-X 'main.version=$(VERSION)' -X 'main.when=$(WHEN)' -X 'main.who=$(WHO)' -X 'main.where=$(WHERE)'"
	strip leatherman
	./leatherman version
