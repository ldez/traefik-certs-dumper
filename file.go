package main

import (
	"encoding/json"
	"os"

	"github.com/fsnotify/fsnotify"
)

// FileBackend stores the config for file backend
type FileBackend struct {
	Name string
	Path string
}

func getStoredDataFromFile(path string) (*StoredData, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	data := &StoredData{}
	if err = json.NewDecoder(f).Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

func (b FileBackend) loop(watch bool) (<-chan *StoredData, <-chan error) {

	dataCh := make(chan *StoredData)
	errCh := make(chan error)
	go func() {
		data, err := getStoredDataFromFile(b.Path)
		if err != nil {
			errCh <- err
		}
		dataCh <- data
		if !watch {
			close(dataCh)
			close(errCh)
		}
	}()

	if !watch {
		return dataCh, errCh
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		errCh <- err
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					data, err1 := getStoredDataFromFile(b.Path)
					if err1 != nil {
						errCh <- err1
					}
					dataCh <- data
				}
			case err1, ok := <-watcher.Errors:
				if !ok {
					return
				}
				errCh <- err1
			}
		}
	}()

	err = watcher.Add(b.Path)
	if err != nil {
		errCh <- err
	}

	return dataCh, errCh
}
