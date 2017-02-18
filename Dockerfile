FROM alpine:latest
ENV GOPATH /opt/go
ENV SRCPATH $GOPATH/src/github.com/myENA/consul-backinator
COPY . $SRCPATH
RUN \
	apk add --no-cache --no-progress ca-certificates bash git glide go musl-dev && \
	mkdir -p $GOPATH/bin && \
	export PATH=$GOPATH/bin:$PATH && \
	cd $SRCPATH && \
	chmod +x build/build.sh && \
	build/build.sh -i && \
	mv consul-backinator /usr/local/bin/consul-backinator && \
	apk del --no-cache --no-progress --purge bash git glide go musl-dev && \
	rm -rf $GOPATH /tmp/* /root/.glide
ENTRYPOINT ["/usr/local/bin/consul-backinator"]
