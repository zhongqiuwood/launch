maDEP := $(shell command -v dep 2> /dev/null)

get_tools:
ifndef DEP
	@echo "Installing dep"
	go get -u -v github.com/golang/dep/cmd/dep
else
	@echo "Dep is already installed..."
endif

get_vendor_deps:
	@echo "--> Generating vendor directory via dep ensure"
	@rm -rf .vendor-new
	@dep ensure -v -vendor-only

update_vendor_deps:
	@echo "--> Running dep ensure"
	@rm -rf .vendor-new
	@dep ensure -v -update

install:
	go install ./cmd/okdexd
	go install ./cmd/okdexcli


build_okdex_bins:
	go build $(BUILD_FLAGS) -o build/okdexd ./cmd/okdexd
	go build $(BUILD_FLAGS) -o build/okdexcli ./cmd/okdexcli

# build: build_okdex_bins
# 	go build $(BUILD_FLAGS) -o build/gaiad ${GOPATH}/src/github.com/cosmos/cosmos-sdk/cmd/gaia/cmd/gaiad
# 	go build $(BUILD_FLAGS) -o build/gaiacli ${GOPATH}/src/github.com/cosmos/cosmos-sdk/cmd/gaia/cmd/gaiacli
# 	#go build $(BUILD_FLAGS) -o build/gaiareplay ${GOPATH}/src/github.com/cosmos/cosmos-sdk/cmd/gaia/cmd/gaiareplay
# 	#go build $(BUILD_FLAGS) -o build/gaiakeyutil ${GOPATH}/src/github.com/cosmos/cosmos-sdk/cmd/gaia/cmd/gaiakeyutil

build_gaia_bins:
	go build $(BUILD_FLAGS) -o build/gaiad ${GOPATH}/src/github.com/cosmos/cosmos-sdk/cmd/gaia/cmd/gaiad
	go build $(BUILD_FLAGS) -o build/gaiacli ${GOPATH}/src/github.com/cosmos/cosmos-sdk/cmd/gaia/cmd/gaiacli
	#go build $(BUILD_FLAGS) -o build/gaiareplay ${GOPATH}/src/github.com/cosmos/cosmos-sdk/cmd/gaia/cmd/gaiareplay
	#go build $(BUILD_FLAGS) -o build/gaiakeyutil ${GOPATH}/src/github.com/cosmos/cosmos-sdk/cmd/gaia/cmd/gaiakeyutil

# build-linux: vendor-deps
# 	LEDGER_ENABLED=false GOOS=linux GOARCH=amd64 $(MAKE) build

# gaiad:
# 	LEDGER_ENABLED=false GOOS=${GOOS} GOARCH=amd64 $(MAKE) build_gaiad

# build_gaiad:
# 	go build -o build/gaiad ${GOPATH}/src/github.com/cosmos/cosmos-sdk/cmd/gaia/cmd/gaiad

# mac:
# 	LEDGER_ENABLED=false GOOS=${GOOS} GOARCH=${GOARCH} $(MAKE) build_okdex_bins

# macgaiad:
# 	LEDGER_ENABLED=false GOOS=${GOOS} GOARCH=${GOARCH} $(MAKE) gaiad



dex-base:
	@./docker/provision/buildbase.sh 0.2

dex-src:
	@echo "Building docker src-image start:  $(@D)"
	mkdir -p $(@D)
	cat docker/src/Dockerfile.in \
		| sed -e 's/_TAG_/$(DOCKER_TAG)/g' \
		> $(@D)/Dockerfile
	@git ls-files | tar -jcT - > $(@D)/okdex.tar.bz2
	docker build -t okchain/dex-src $(@D)

dex-bin: docker_okdexd docker_okdexcli

dex-node: dex-src dex-bin
	@echo "Building docker node-image start:  $(@D)"
	mkdir -p $(@D)
	cat docker/node/Dockerfile.in \
		| sed -e 's/_TAG_/$(DOCKER_TAG)/g' \
		> $(@D)/Dockerfile
	docker build -t okchain/dex-node $(@D)


docker_okdexd:
	@echo "Building $@"
	@docker run -i \
    		--user=$(UID) \
    		-v $(abspath build):/opt/gopath/bin \
    		okchain/dex-src:latest go install -ldflags "$(GO_LDFLAGS)" github.com/okchain/okdex/cmd/okdexd

docker_okdexcli:
	@echo "Building $@"
	@docker run -i \
    		--user=$(UID) \
    		-v $(abspath build):/opt/gopath/bin \
    		okchain/dex-src:latest go install -ldflags "$(GO_LDFLAGS)" github.com/okchain/okdex/cmd/okdexcli

format:
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/lcd/statik/statik.go" | xargs gofmt -w -s


.PHONY: build