GOOPTS := GOARCH=amd64 CGO_ENABLED=0

build: linux windows macos

linux: modules bin generate
	$(GOOPTS) GOOS=linux go build -o bin/packer-post-processor-teamcity.linux

windows: modules bin generate
	$(GOOPTS) GOOS=windows go build -o bin/packer-post-processor-teamcity.exe

macos: modules bin generate
	$(GOOPTS) GOOS=darwin go build -o bin/packer-post-processor-teamcity.macos

modules:
	go mod download

tools:
	go install github.com/hashicorp/packer/cmd/mapstructure-to-hcl2

generate: tools
	go generate ./...

bin:
	mkdir -p bin
	rm -f bin/*

.PHONY: bin
