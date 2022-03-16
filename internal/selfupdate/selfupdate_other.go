//go:build !linux
// +build !linux

package selfupdate

func isSameFile(string) error { return nil }
