module github.com/ldez/traefik-certs-dumper/v2

go 1.14

require (
	github.com/abronan/valkeyrie v0.0.0-20200127174252-ef4277a138cd
	github.com/fsnotify/fsnotify v1.4.8-0.20190312181446-1485a34d5d57
	github.com/go-acme/lego/v4 v4.0.1
	github.com/mitchellh/go-homedir v1.1.0
	github.com/spf13/cobra v0.0.7
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.6.1
	github.com/traefik/traefik/v2 v2.3.1
)

// related to Traefik
replace github.com/docker/docker => github.com/docker/engine v0.0.0-20190725163905-fa8dd90ceb7b

// related to Traefik: Containous forks
replace (
	github.com/abbot/go-http-auth => github.com/containous/go-http-auth v0.4.1-0.20180112153951-65b0cdae8d7f
	github.com/go-check/check => github.com/containous/check v0.0.0-20170915194414-ca0bf163426a
	github.com/gorilla/mux => github.com/containous/mux v0.0.0-20181024131434-c33f32e26898
	github.com/mailgun/minheap => github.com/containous/minheap v0.0.0-20190809180810-6e71eb837595
	github.com/mailgun/multibuf => github.com/containous/multibuf v0.0.0-20190809014333-8b6c9a7e6bba
	github.com/rancher/go-rancher-metadata => github.com/containous/go-rancher-metadata v0.0.0-20190402144056-c6a65f8b7a28
)
