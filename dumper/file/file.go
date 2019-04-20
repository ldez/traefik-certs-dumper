package file

import (
	"encoding/json"
	"os"

	"github.com/ldez/traefik-certs-dumper/dumper"
)

func Dump(acmeFile, dumpPath string, crtInfo, keyInfo dumper.FileInfo, domainSubDir bool) error {
	data, err := readFile(acmeFile)
	if err != nil {
		return err
	}

	return dumper.Dump(data, dumpPath, crtInfo, keyInfo, domainSubDir)
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
