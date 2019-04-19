FROM golang:1-alpine as builder

RUN apk --update upgrade \
    && apk --no-cache --no-progress add git make gcc musl-dev

WORKDIR /go/src/github.com/ldez/traefik-certs-dumper

ENV GO111MODULE on
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN make build

FROM alpine:3.9
RUN apk --update upgrade \
    && apk --no-cache --no-progress add ca-certificates \
    && update-ca-certificates

COPY --from=builder /go/src/github.com/ldez/traefik-certs-dumper/traefik-certs-dumper /usr/bin/traefik-certs-dumper

ENTRYPOINT ["/usr/bin/traefik-certs-dumper"]
