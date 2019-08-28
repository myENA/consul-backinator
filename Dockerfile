FROM golang:alpine as build
ENV BLDPATH /opt/consul-backinator
COPY . $BLDPATH
RUN \
	apk add --no-cache --no-progress ca-certificates bash git musl-dev && \
	cd $BLDPATH && \
	chmod +x build/build.sh && \
	build/build.sh && \
	mv consul-backinator /usr/local/bin/consul-backinator

FROM alpine:latest
COPY --from=build /usr/local/bin/consul-backinator /usr/local/bin/consul-backinator
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ENTRYPOINT ["/usr/local/bin/consul-backinator"]
