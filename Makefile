VERSION=$(shell git describe --tags --candidates=1 --dirty)
BUILD_FLAGS=-ldflags="-X main.Version=$(VERSION)" -trimpath
SRC=$(shell find . -name '*.go') go.mod
INSTALL_DIR ?= ~/bin
.PHONY: binaries clean release install

aws-vault: $(SRC)
	go build -ldflags="-X main.Version=$(VERSION)" -o $@ .

install: aws-vault
	mkdir -p $(INSTALL_DIR)
	rm -f $(INSTALL_DIR)/aws-vault
	cp -a ./aws-vault $(INSTALL_DIR)/aws-vault

binaries: aws-vault-linux-amd64 aws-vault-linux-arm64 aws-vault-linux-ppc64le aws-vault-linux-arm7 aws-vault-darwin-amd64 aws-vault-darwin-arm64 aws-vault-windows-386.exe aws-vault-windows-arm64.exe aws-vault-freebsd-amd64

clean:
	rm -f ./aws-vault ./aws-vault-*-* ./SHA256SUMS

release: binaries SHA256SUMS

	@echo "\nTo create a new release run:\n\n    gh release create --title $(VERSION) $(VERSION) \
	aws-vault-freebsd-amd64 \
	aws-vault-linux-amd64 \
	aws-vault-linux-arm64 \
	aws-vault-linux-arm7 \
	aws-vault-linux-ppc64le \
	aws-vault-darwin-amd64 \
	aws-vault-darwin-arm64 \
	aws-vault-windows-386.exe \
	aws-vault-windows-arm64.exe \
	SHA256SUMS\n"

aws-vault-darwin-amd64: $(SRC)
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 SDKROOT=$(shell xcrun --sdk macosx --show-sdk-path) go build $(BUILD_FLAGS) -o $@ .

aws-vault-darwin-arm64: $(SRC)
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 SDKROOT=$(shell xcrun --sdk macosx --show-sdk-path) go build $(BUILD_FLAGS) -o $@ .

aws-vault-freebsd-amd64: $(SRC)
	GOOS=freebsd GOARCH=amd64 go build $(BUILD_FLAGS) -o $@ .

aws-vault-linux-amd64: $(SRC)
	GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o $@ .

aws-vault-linux-arm64: $(SRC)
	GOOS=linux GOARCH=arm64 go build $(BUILD_FLAGS) -o $@ .

aws-vault-linux-ppc64le: $(SRC)
	GOOS=linux GOARCH=ppc64le go build $(BUILD_FLAGS) -o $@ .

aws-vault-linux-arm7: $(SRC)
	GOOS=linux GOARCH=arm GOARM=7 go build $(BUILD_FLAGS) -o $@ .

aws-vault-windows-386.exe: $(SRC)
	GOOS=windows GOARCH=386 go build $(BUILD_FLAGS) -o $@ .

aws-vault-windows-arm64.exe: $(SRC)
	GOOS=windows GOARCH=arm64 go build $(BUILD_FLAGS) -o $@ .

SHA256SUMS: binaries
	shasum -a 256 \
	  aws-vault-freebsd-amd64 \
	  aws-vault-linux-amd64 \
	  aws-vault-linux-arm64 \
	  aws-vault-linux-arm7 \
	  aws-vault-linux-ppc64le \
	  aws-vault-darwin-amd64 \
	  aws-vault-darwin-arm64 \
	  aws-vault-windows-386.exe \
	  aws-vault-windows-arm64.exe \
	    > $@
