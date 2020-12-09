FROM golang:alpine AS build-stage
WORKDIR /go/src/consul-backinator
COPY . .
RUN apk add --no-cache --no-progress ca-certificates tzdata git musl-dev && \
    go install



FROM alpine:latest

WORKDIR /usr/local/bin
COPY --from=build-stage /go/bin/consul-backinator ./
RUN \
	apk add --no-cache --no-progress ca-certificates tzdata
ENTRYPOINT ["/usr/local/bin/consul-backinator"]
