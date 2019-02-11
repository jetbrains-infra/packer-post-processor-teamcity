GOOPTS := GOARCH=amd64 CGO_ENABLED=0

build: linux windows macos

linux: modules bin
	$(GOOPTS) GOOS=linux go build -o bin/packer-post-processor-teamcity.linux

windows: modules bin
	$(GOOPTS) GOOS=windows go build -o bin/packer-post-processor-teamcity.exe

macos: modules bin
	$(GOOPTS) GOOS=darwin go build -o bin/packer-post-processor-teamcity.macos

modules:
	go mod download

bin:
	mkdir -p bin
	rm -f bin/*

.PHONY: bin
