IMAGE := till/law:$(VERSION)

.DEFAULT_GOAL := build

build:
	go build *.go

test:
	go test $$(go list ./... | grep -v vendor)

image: require-version
	docker build -t $(IMAGE) .

require-version:
ifndef VERSION
	$(error VERSION needs to be set)
endif