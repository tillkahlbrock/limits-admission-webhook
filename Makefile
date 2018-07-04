.DEFAULT_GOAL := build

build:
	go build *.go

test:
	go test $$(go list ./... | grep -v vendor)
