VERSION := $(shell git describe --abbrev=7 --dirty --always || echo $TRAVIS_COMMIT)
WHEN := $(shell date)
WHO := $(shell whoami)
WHERE := $(shell hostname -f)

cmd/leatherman/leatherman.xz: cmd/leatherman/leatherman
	xz cmd/leatherman/leatherman

cmd/leatherman/leatherman:
	( cd / ; go get -u github.com/golang/lint/golint )
	export GO111MODULE=on
	go get -t ./...
	golint -set_exit_status ./...
	go vet ./...
	TZ=America/Los_Angeles go test ./...
	( cd cmd/leatherman; go build -ldflags "-s -X 'main.version=$(VERSION)' -X 'main.when=$(WHEN)' -X 'main.who=$(WHO)' -X 'main.where=$(WHERE)'" )
	cmd/leatherman/leatherman version
