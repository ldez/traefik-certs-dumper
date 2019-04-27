FROM golang:1-alpine as builder

ARG RUNTIME_HASH
ARG GOARCH
ARG GOARM
ARG GOOS

RUN apk --update upgrade \
    && apk --no-cache --no-progress add git make gcc musl-dev

WORKDIR /go/src/github.com/ldez/traefik-certs-dumper

ENV GO111MODULE on
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN GOARCH=${GOARCH} GOARM=${GOARM} GOOS=${GOOS} make build

FROM alpine:3.9${RUNTIME_HASH}
RUN apk --update upgrade \
    && apk --no-cache --no-progress add ca-certificates

COPY --from=builder /go/src/github.com/ldez/traefik-certs-dumper/traefik-certs-dumper /usr/bin/traefik-certs-dumper

ENTRYPOINT ["/usr/bin/traefik-certs-dumper"]
