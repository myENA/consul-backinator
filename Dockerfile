FROM alpine:latest
ENV GOPATH /opt/go
ENV SRCPATH $GOPATH/src/github.com/myENA/consul-backinator
COPY . $SRCPATH
RUN \
	apk add --no-cache --no-progress ca-certificates go git bash && \
	mkdir -p $GOPATH/bin && \
	export PATH=$GOPATH/bin:$PATH && \
	cd $SRCPATH && \
	chmod +x build/build.sh && \
	build/build.sh -i && \
	mv consul-backinator /usr/local/bin/consul-backinator && \
	apk del --no-cache --no-progress --purge go git bash && \
	rm -rf $GOPATH /tmp/* /root/.glide
ENTRYPOINT ["/usr/local/bin/consul-backinator"]
