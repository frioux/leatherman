VERSION := $(shell git describe --abbrev=7 --dirty --always || echo $TRAVIS_COMMIT)
WHEN := $(shell date)
WHO := $(shell whoami)
WHERE := $(shell hostname -f)

leatherman.xz: leatherman
	xz --stdout leatherman > leatherman.xz

leatherman: export GO111MODULE = on
leatherman:
	GO111MODULE=off go get -u golang.org/x/lint/golint
	go get ./...
	golint -set_exit_status ./...
	go vet ./...
	TZ=America/Los_Angeles go test -coverprofile=cover.cover -race ./...
	go mod verify
	go build -ldflags "-s -X 'main.version=$(VERSION)' -X 'main.when=$(WHEN)' -X 'main.who=$(WHO)' -X 'main.where=$(WHERE)'"
	./leatherman version

watch:
	minotaur . -- ./internal/build-test
