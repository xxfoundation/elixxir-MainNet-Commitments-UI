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

linux64_binary:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags '-w -s' -o mainnet-commitments-ui.linux64 main.go index.go

win64_binary:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags '-w -s' -o mainnet-commitments-ui.win64 main.go index.go

darwin64_binary:
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags '-w -s' -o mainnet-commitments-ui.darwin64 main.go index.go


binaries: linux64_binary win64_binary darwin64_binary