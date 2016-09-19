SUDO ?=

ifeq ($(SUDO),true)
	sudo = sudo
endif

.PHONY: build release docker

build:
	@build/build.sh -i

release:
	@build/build.sh -ri

docker: release
	$(sudo) docker build -t my-ena/consul-backinator -f build/docker .
