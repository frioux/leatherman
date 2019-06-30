VERSION := $(shell git describe --abbrev=7 --dirty --always)
WHEN := $(shell git log -1 --pretty=%cI $(VERSION) 2>/dev/null)

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
	go build -ldflags "-s -X 'github.com/frioux/leatherman/internal/version.Version=$(VERSION)' -X 'github.com/frioux/leatherman/internal/version.When=$(WHEN)'"
	./leatherman version

watch:
	minotaur . -- ./internal/build-test
