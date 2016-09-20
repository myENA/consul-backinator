REGISTRY ?=
IMAGE_PATH ?= $(REGISTRY)ena/consul-backinator
SUDO ?=
RELEASE_VERSION = $(shell grep RELEASE_VERSION= build/build.sh | grep -oE '[0-9]+?\.[0-9]+?')

ifeq ($(SUDO),true)
	sudo = sudo
endif

.PHONY: build release check clean docker docker_release

build:
	@build/build.sh -i

release:
	@build/build.sh -ir

check:
	@build/codeCheck.sh

clean:
	@build/build.sh -d

docker: release
	$(sudo) docker build -t $(IMAGE_PATH):$(RELEASE_VERSION) -f build/docker . && \
	$(sudo) docker tag $(IMAGE_PATH):$(RELEASE_VERSION) $(IMAGE_PATH):latest

docker_release: docker
	$(sudo) docker push $(IMAGE_PATH):$(RELEASE_VERSION) && \
	$(sudo) docker push $(IMAGE_PATH):latest
