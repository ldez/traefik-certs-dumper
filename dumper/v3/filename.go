//go:build !windows
// +build !windows

package v3

func safeName(filename string) string {
	return filename
}
