package file

import (
	"testing"

	"github.com/ldez/traefik-certs-dumper/v2/dumper"
	"github.com/stretchr/testify/require"
)

func TestDump(t *testing.T) {
	testCases := []struct {
		desc     string
		acmeFile string
		version  string
	}{
		{
			desc:     "should skip EOF error",
			acmeFile: "./fixtures/acme-empty.json",
		},
		{
			desc:     "should dump traefik v1 file content",
			acmeFile: "./fixtures/acme-v1.json",
		},
		{
			desc:     "should dump traefik v2 file content",
			acmeFile: "./fixtures/acme-v2.json",
			version:  "v2",
		},
		{
			desc:     "should dump traefik v3 file content",
			acmeFile: "./fixtures/acme-v3.json",
			version:  "v3",
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			dir := t.TempDir()

			cfg := &dumper.BaseConfig{
				DumpPath: dir,
				CrtInfo: dumper.FileInfo{
					Name: "certificate",
					Ext:  ".crt",
				},
				KeyInfo: dumper.FileInfo{
					Name: "privatekey",
					Ext:  ".key",
				},
				Clean:   true,
				Version: test.version,
			}

			err := Dump(t.Context(), test.acmeFile, cfg)
			require.NoError(t, err)
		})
	}
}
