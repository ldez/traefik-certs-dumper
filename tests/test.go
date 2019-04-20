package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/abronan/valkeyrie"
	"github.com/abronan/valkeyrie/store"
	"github.com/abronan/valkeyrie/store/boltdb"
	"github.com/abronan/valkeyrie/store/consul"
	etcdv3 "github.com/abronan/valkeyrie/store/etcd/v3"
	"github.com/abronan/valkeyrie/store/zookeeper"
)

const (
	storeKey = "traefik/acme/account/object"
)

func main() {
	writeDataToBackends()
}

func writeDataToBackends() {
	storeConfig := &store.Config{
		ConnectionTimeout: 3 * time.Second,
		Bucket:            "traefik",
	}

	consul.Register()
	etcdv3.Register()
	zookeeper.Register()
	boltdb.Register()

	consulStore, err := valkeyrie.NewStore(
		store.CONSUL,
		[]string{"localhost:8500"},
		storeConfig,
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	etcdv3Store, err := valkeyrie.NewStore(
		store.ETCDV3,
		[]string{"localhost:2379"},
		storeConfig,
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	zkStore, err := valkeyrie.NewStore(
		store.ZK,
		[]string{"localhost:2181"},
		storeConfig,
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	boltdbStore, err := valkeyrie.NewStore(
		store.BOLTDB,
		[]string{"/tmp/my.db"},
		storeConfig,
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	f, _ := os.Open("/tmp/acme.json")
	reader := bufio.NewReader(f)
	content, _ := ioutil.ReadAll(reader)

	var b bytes.Buffer
	gz := gzip.NewWriter(&b)

	_, err = gz.Write(content)
	if err != nil {
		return
	}

	if err = gz.Flush(); err != nil {
		return
	}

	if err = gz.Close(); err != nil {
		return
	}

	if err := boltdbStore.Put(storeKey, b.Bytes(), nil); err == nil {
		fmt.Println("successfully updated boltdb")
	}
	if err := zkStore.Put(storeKey, b.Bytes(), nil); err == nil {
		fmt.Println("successfully updated zookeeper")
	}
	if err := etcdv3Store.Put(storeKey, b.Bytes(), nil); err == nil {
		fmt.Println("successfully updated etcd")
	}
	if err := consulStore.Put(storeKey, b.Bytes(), nil); err == nil {
		fmt.Println("successfully updated consul")
	}
}
