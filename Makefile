VERSION := $(shell git describe --abbrev=7 --dirty --always || echo $TRAVIS_COMMIT)
WHEN := $(shell date)
WHO := $(shell whoami)
WHERE := $(shell hostname -f)

leatherman.xz: leatherman
	xz leatherman

leatherman: *.go
	( cd / ; go get -u github.com/golang/lint/golint )
	export GO111MODULE=on
	go get -t ./...
	golint -set_exit_status ./...
	go vet ./...
	TZ=America/Los_Angeles go test ./...
	go build -ldflags "-s -X 'main.version=$(VERSION)' -X 'main.when=$(WHEN)' -X 'main.who=$(WHO)' -X 'main.where=$(WHERE)'"
	./leatherman version
