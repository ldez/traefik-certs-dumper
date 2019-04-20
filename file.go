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

func (b FileBackend) getStoredData(watch bool) (<-chan *StoredData, <-chan error) {
	dataCh := make(chan *StoredData)
	errCh := make(chan error)
	go func() {
		sendStoredData(b.Path, dataCh, errCh)
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

	go loopFile(b.Path, watcher, dataCh, errCh)

	err = watcher.Add(b.Path)
	if err != nil {
		errCh <- err
	}

	return dataCh, errCh
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

func sendStoredData(path string, dataCh chan *StoredData, errCh chan error) {
	data, err := getStoredDataFromFile(path)
	if err != nil {
		errCh <- err
	}
	dataCh <- data
}

func loopFile(path string, watcher *fsnotify.Watcher, dataCh chan *StoredData, errCh chan error) {
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					sendStoredData(path, dataCh, errCh)
				}
			case err1, ok := <-watcher.Errors:
				if !ok {
					return
				}
				errCh <- err1
			}
		}
	}()
}
