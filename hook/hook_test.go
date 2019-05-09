package hook

import "testing"

func Test_execute(t *testing.T) {
	testCases := []struct {
		desc    string
		command string
	}{
		{
			desc:    "expand env vars",
			command: `echo "${GOPATH} ${GOARCH}"`,
		},
		{
			desc:    "simple",
			command: `echo 'hello'`,
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {

			err := execute(test.command)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
