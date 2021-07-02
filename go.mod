module github.com/ldez/traefik-certs-dumper/v2

go 1.16

require (
	github.com/abronan/valkeyrie v0.1.0
	github.com/fsnotify/fsnotify v1.4.9
	github.com/go-acme/lego/v4 v4.4.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/spf13/cobra v1.2.0
	github.com/spf13/viper v1.8.1
	github.com/stretchr/testify v1.7.0
	github.com/traefik/traefik/v2 v2.4.9
)

// Containous forks
replace (
	github.com/abbot/go-http-auth => github.com/containous/go-http-auth v0.4.1-0.20200324110947-a37a7636d23e
	github.com/go-check/check => github.com/containous/check v0.0.0-20170915194414-ca0bf163426a
	github.com/gorilla/mux => github.com/containous/mux v0.0.0-20181024131434-c33f32e26898
	github.com/mailgun/minheap => github.com/containous/minheap v0.0.0-20190809180810-6e71eb837595
	github.com/mailgun/multibuf => github.com/containous/multibuf v0.0.0-20190809014333-8b6c9a7e6bba
)
