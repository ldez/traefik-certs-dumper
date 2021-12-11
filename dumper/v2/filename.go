//go:build !windows
// +build !windows

package v2

func safeName(filename string) string {
	return filename
}
