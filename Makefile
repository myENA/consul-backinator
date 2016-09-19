DEPS := $(shell git ls-files '*.go' | grep -v '^vendor')
SUDO ?=

ifeq ($(SUDO),true)
	sudo = sudo
endif

.phony: binary bootstrap docker

binary: bootstrap consul-backinator

bootstrap:
	go get -d

consul-backinator: $(DEPS)
	CGO_ENABLED=0 go build -o $@ -x -a -installsuffix cgo -ldflags '-s'

docker: binary
	$(sudo) docker build -t my-ena/consul-backinator .
