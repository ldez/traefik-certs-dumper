package file

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"io"
	"log"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/ldez/traefik-certs-dumper/v2/dumper"
	"github.com/ldez/traefik-certs-dumper/v2/hook"
)

// Dump Dumps "acme.json" file to certificates.
func Dump(acmeFile string, baseConfig *dumper.BaseConfig) error {
	err := dump(acmeFile, baseConfig)
	if err != nil {
		return err
	}

	if baseConfig.Watch {
		return watch(acmeFile, baseConfig)
	}
	return nil
}

func dump(acmeFile string, baseConfig *dumper.BaseConfig) error {
	data, err := readFile(acmeFile)
	if err != nil {
		return err
	}

	return dumper.Dump(data, baseConfig)
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

func watch(acmeFile string, baseConfig *dumper.BaseConfig) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	defer func() { _ = watcher.Close() }()

	done := make(chan bool)
	go func() {
		var previousHash []byte

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if isDebug() {
					log.Println("event:", event)
				}

				hash, errW := manageEvent(watcher, event, acmeFile, previousHash, baseConfig)
				if errW != nil {
					log.Println("error:", errW)
					done <- true
					return
				}

				previousHash = hash

			case errW, ok := <-watcher.Errors:
				if !ok {
					return
				}

				log.Println("error:", errW)
				done <- true
				return
			}
		}
	}()

	err = watcher.Add(acmeFile)
	if err != nil {
		return err
	}

	<-done

	return nil
}

func manageEvent(watcher *fsnotify.Watcher, event fsnotify.Event, acmeFile string, previousHash []byte, baseConfig *dumper.BaseConfig) ([]byte, error) {
	err := manageRename(watcher, event, acmeFile)
	if err != nil {
		return nil, err
	}

	hash, err := calculateHash(acmeFile)
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(previousHash, hash) {
		if isDebug() {
			log.Println("detected changes on file:", event.Name)
		}

		if errD := dump(acmeFile, baseConfig); errD != nil {
			return nil, errD
		}

		if isDebug() {
			log.Println("Dumped new certificate data.")
		}

		hook.Exec(baseConfig.Hook)
	}

	return hash, nil
}

func manageRename(watcher *fsnotify.Watcher, event fsnotify.Event, acmeFile string) error {
	if event.Op&fsnotify.Rename != fsnotify.Rename {
		return nil
	}

	if err := watcher.Remove(acmeFile); err != nil {
		return err
	}

	return watcher.Add(acmeFile)
}

func calculateHash(acmeFile string) ([]byte, error) {
	file, err := os.Open(acmeFile)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	h := md5.New()
	_, err = io.Copy(h, file)
	if err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

func isDebug() bool {
	return strings.EqualFold(os.Getenv("TCD_DEBUG"), "true")
}
