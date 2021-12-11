//go:build !windows
// +build !windows

package v1

func safeName(filename string) string {
	return filename
}
