.PHONY: update master release setup update_master update_release build clean

clean:
	rm -rf vendor/
	go mod vendor

update:
	-GOFLAGS="" go get all

build:
	go build ./...
	go mod tidy

update_master:
	GOFLAGS="" go get -d git.xx.network/elixxir/mainnet-commitments@jonah/tm-change

master: update_master clean build
