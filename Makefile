SUDO ?=

ifeq ($(SUDO),true)
	sudo = sudo
endif

.PHONY: build release clean docker

build:
	@build/build.sh -i

release:
	@build/build.sh -ir

clean:
	@build/build.sh -d

docker: release
	$(sudo) docker build -t ena/consul-backinator -f build/docker .
