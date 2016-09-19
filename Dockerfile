FROM scratch
COPY ca-certificates.crt /etc/ssl/certs/
COPY consul-backinator /
ENTRYPOINT ["/consul-backinator"]
