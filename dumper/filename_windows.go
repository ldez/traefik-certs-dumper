// +build windows

package dumper

import "strings"

func safeName(filename string) string {
	return strings.ReplaceAll(filename, "*", "_")
}
