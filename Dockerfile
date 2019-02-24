FROM golang:1-alpine as builder

RUN apk --update upgrade \
&& apk --no-cache --no-progress add git make gcc musl-dev \
&& rm -rf /var/cache/apk/*

WORKDIR /go/src/github.com/ldez/traefik-certs-dumper
COPY . .

RUN go get -u github.com/golang/dep/cmd/dep
ENV GO111MODULE on
RUN go mod download
RUN make build

FROM alpine:3.9
RUN apk --update upgrade \
    && apk --no-cache --no-progress add ca-certificates git \
    && rm -rf /var/cache/apk/*

COPY --from=builder /go/src/github.com/ldez/traefik-certs-dumper/traefik-certs-dumper /usr/bin/traefik-certs-dumper

ENTRYPOINT ["/usr/bin/traefik-certs-dumper"]
