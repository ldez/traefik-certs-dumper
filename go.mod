module github.com/ldez/traefik-certs-dumper/v2

go 1.12

require (
	github.com/abronan/valkeyrie v0.0.0-20190419181538-ccf7df650fe4
	github.com/cenkalti/backoff v2.1.1+incompatible // indirect
	github.com/fsnotify/fsnotify v1.4.7
	github.com/go-acme/lego v2.7.2+incompatible
	github.com/hashicorp/go-msgpack v0.5.4 // indirect
	github.com/hashicorp/go-uuid v1.0.1 // indirect
	github.com/hashicorp/memberlist v0.1.3 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/go-testing-interface v1.0.0 // indirect
	github.com/pascaldekloe/goe v0.1.0 // indirect
	github.com/pkg/errors v0.8.1 // indirect
	github.com/sirupsen/logrus v1.4.1 // indirect
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.4.0
	gopkg.in/square/go-jose.v2 v2.3.1 // indirect
)

replace (
	github.com/ugorji/go => github.com/ugorji/go v1.1.2-0.20181022190402-e5e69e061d4f
	github.com/ugorji/go/codec v0.0.0-20181204163529-d75b2dcb6bc8 => github.com/ugorji/go/codec v1.1.2-0.20181022190402-e5e69e061d4f
)
