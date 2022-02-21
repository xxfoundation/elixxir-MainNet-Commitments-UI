.PHONY: update master release setup update_master update_release build clean

clean:
	rm -rf vendor/
	go mod vendor

update:
	-GOFLAGS="" go get all

build:
	go build ./...
	go mod tidy

update_release:
	GOFLAGS="" go get gitlab.com/elixxir/client@release

update_master:
	GOFLAGS="" go get gitlab.com/elixxir/client@master

master: update_master clean build

release: update_release clean build
