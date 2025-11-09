# syntax=docker/dockerfile:1.4
FROM alpine:3

ARG TARGETPLATFORM

RUN apk --no-cache --no-progress add git ca-certificates tzdata jq \
    && rm -rf /var/cache/apk/*

COPY $TARGETPLATFORM/traefik-certs-dumper /usr/bin/traefik-certs-dumper

ENTRYPOINT ["/usr/bin/traefik-certs-dumper"]
