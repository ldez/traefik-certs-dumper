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

				if strings.EqualFold(os.Getenv("TCD_DEBUG"), "true") {
					log.Println("event:", event)
				}

				switch {
				case event.Op&fsnotify.Rename == fsnotify.Rename:
					errW := watcher.Remove(acmeFile)
					if errW != nil {
						log.Println("error:", errW)
						done <- true
						return
					}

					errW = watcher.Add(acmeFile)
					if errW != nil {
						log.Println("error:", errW)
						done <- true
						return
					}
					fallthrough
				case event.Op&fsnotify.Write == fsnotify.Write:
					hash, errH := calculateHash(acmeFile)
					if err != nil {
						log.Println("error:", errH)
						done <- true
						return
					}

					if !bytes.Equal(previousHash, hash) {
						previousHash = hash

						if strings.EqualFold(os.Getenv("TCD_DEBUG"), "true") {
							log.Println("detected changes on file:", event.Name)
						}

						if errD := dump(acmeFile, baseConfig); errD != nil {
							log.Println("error:", errD)
							done <- true
							return
						}

						log.Println("Dumped new certificate data.")
					}

				}
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
