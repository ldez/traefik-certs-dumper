package file

import (
	"encoding/json"
	"log"
	"os"

	"github.com/ldez/traefik-certs-dumper/v2/dumper"
	"github.com/rjeczalik/notify"
)

// Dump Dumps "acme.json" file to certificates.
func Dump(config *Config, baseConfig *dumper.BaseConfig) error {
	if config.Watch {
		return watch(config.AcmeFile, baseConfig)
	}

	data, err := readFile(config.AcmeFile)
	if err != nil {
		return err
	}

	return dumper.Dump(data, baseConfig)
}

func watch(acmeFile string, baseConfig *dumper.BaseConfig) error {
	events := make(chan notify.EventInfo, 1)

	if err := notify.Watch(acmeFile, events, notify.Create, notify.Write, notify.Remove); err != nil {
		return err
	}
	defer notify.Stop(events)

	for {
		// wait for filesystem event
		<-events

		data, err := readFile(acmeFile)
		if err != nil {
			return err
		}

		if err := dumper.Dump(data, baseConfig); err != nil {
			return err
		}

		log.Println("Dumped new certificate data.")
	}
}

func readFile(acmeFile string) (*dumper.StoredData, error) {
	source, err := os.Open(acmeFile)
	if err != nil {
		return nil, err
	}

	data := &dumper.StoredData{}
	if err = json.NewDecoder(source).Decode(data); err != nil {
		return nil, err
	}

	return data, nil
}
