//go:build windows
// +build windows

package v3

import "strings"

func safeName(filename string) string {
	return strings.ReplaceAll(filename, "*", "_")
}
