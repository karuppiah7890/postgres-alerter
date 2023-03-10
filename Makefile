# TODO: Build using https://goreleaser.com/
build:
	CGO_ENABLED=0 go build -v

build-linux:
	GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -v

release-linux: build-linux
	tar cvzf postgres-alerter-linux-amd64.tar.gz postgres-alerter

# TODO: Lint using golangci-lint
