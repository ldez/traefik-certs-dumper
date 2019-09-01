// +build windows

package v2

import "strings"

func safeName(filename string) string {
	return strings.ReplaceAll(filename, "*", "_")
}
