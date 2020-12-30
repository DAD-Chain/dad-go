GOFMT=gofmt
GC=go build
VERSION := $(shell git describe --abbrev=4 --dirty --always --tags)
BUILD_NODE_PAR = -ldflags "-X github.com/ontio/dad-go/common/config.Version=$(VERSION)" #-race

ARCH=$(shell uname -m)
DBUILD=docker build
DRUN=docker run
DOCKER_NS ?= ontio
DOCKER_TAG=$(ARCH)-$(VERSION)
ONT_CFG_IN_DOCKER=config-solo.json
WALLET_FILE=wallet.dat

SRC_FILES = $(shell git ls-files | grep -e .go$ | grep -v _test.go)
TOOLS=./tools
NATIVE_ABI=$(TOOLS)/abi/native

dad-go: $(SRC_FILES)
	$(GC)  $(BUILD_NODE_PAR) -o dad-go main.go

tools: $(SRC_FILES)
	$(GC)  $(BUILD_NODE_PAR) -o sigsvr sigsvr.go
	@if [ ! -d $(TOOLS) ];then mkdir $(TOOLS) ;fi
	@mv sigsvr $(TOOLS)
	@if [ ! -d $(NATIVE_ABI) ];then mkdir -p $(NATIVE_ABI) ;fi
	@cp ./cmd/abi/native/*.json $(NATIVE_ABI)

all: dad-go tools
	
format:
	$(GOFMT) -w main.go

$(WALLET_FILE):
	@if [ ! -e $(WALLET_FILE) ]; then $(error Please create wallet file first) ; fi

docker/payload: docker/build/bin/dad-go docker/Dockerfile $(ONT_CFG_IN_DOCKER) $(WALLET_FILE)
	@echo "Building dad-go payload"
	@mkdir -p $@
	@cp docker/Dockerfile $@
	@cp docker/build/bin/dad-go $@
	@cp -f $(ONT_CFG_IN_DOCKER) $@/config.json
	@cp -f $(WALLET_FILE) $@
	@tar czf $@/config.tgz -C $@ config.json $(WALLET_FILE)
	@touch $@

docker/build/bin/%: Makefile
	@echo "Building dad-go in docker"
	@mkdir -p docker/build/bin docker/build/pkg
	@$(DRUN) --rm \
		-v $(abspath docker/build/bin):/go/bin \
		-v $(abspath docker/build/pkg):/go/pkg \
		-v $(GOPATH)/src:/go/src \
		-w /go/src/github.com/ontio/dad-go \
		golang:1.9.5-stretch \
		$(GC)  $(BUILD_NODE_PAR) -o docker/build/bin/dad-go main.go
	@touch $@

docker: Makefile docker/payload docker/Dockerfile 
	@echo "Building dad-go docker"
	@$(DBUILD) -t $(DOCKER_NS)/dad-go docker/payload
	@docker tag $(DOCKER_NS)/dad-go $(DOCKER_NS)/dad-go:$(DOCKER_TAG)
	@touch $@

clean:
	rm -rf *.8 *.o *.out *.6
	rm -rf dad-go tools docker/payload docker/build

