// +build windows

package v1

import "strings"

func safeName(filename string) string {
	return strings.ReplaceAll(filename, "*", "_")
}
