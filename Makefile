REGISTRY ?=
IMAGE_PATH ?= $(REGISTRY)ena/consul-backinator
SUDO ?=
RELEASE_VERSION = $(shell grep RELEASE_VERSION= build/build.sh | grep -oE '[0-9]+?\.[0-9]+?')

ifeq ($(SUDO),true)
	sudo = sudo
endif

.PHONY: build test release check clean distclean docker docker-release

build:
	@build/build.sh -i

test:
	@$(GOPATH)/bin/porter go test -v

release:
	@build/build.sh -ir

check:
	@build/codeCheck.sh

clean:
	@build/build.sh -d

distclean:
	@build/build.sh -dc

docker:
	$(sudo) docker build -t $(IMAGE_PATH):latest .

docker-release:
	$(sudo) docker build -t $(IMAGE_PATH):latest .
	$(sudo) docker tag $(IMAGE_PATH):latest $(IMAGE_PATH):$(RELEASE_VERSION)
