.PHONY: \
build-darwin_arm64 build-darwin_amd64 build-linux_amd64 build-windows_amd64 \
setup-macos_m1 setup-macos_intel setup-linux

.SILENT:

build-darwin_arm64:
	GOOS=darwin GOARCH=arm64 go build -o build/bin/docsel-darwin_arm64 cmd/main.go

build-darwin_amd64:
	GOOS=darwin GOARCH=amd64 go build -o build/bin/docsel-darwin_amd64 cmd/main.go

build-linux_amd64:
	GOOS=linux GOARCH=amd64 go build -o build/bin/docsel-linux_amd64 cmd/main.go

build-windows_amd64:
	GOOS=windows GOARCH=amd64 go build -o build/bin/docsel-windows_amd64.exe cmd/main.go

setup-macos_m1:
	sudo mkdir -p /usr/local/bin && sudo cp build/bin/docsel-darwin_arm64 /usr/local/bin/docsel

setup-macos_intel:
	sudo mkdir -p /usr/local/bin && sudo cp build/bin/docsel-darwin_amd64 /usr/local/bin/docsel

setup-linux:
	sudo mkdir -p /usr/local/bin && sudo cp build/bin/docsel-linux_amd64 /usr/local/bin/docsel
