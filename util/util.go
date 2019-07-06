package util

import (
	"os"
	"path/filepath"
	"runtime"
)

// GetCurrentDirectory returns the current path.
func GetCurrentDirectory() string {
	if runtime.GOOS == "windows" {
		path, err := os.Getwd()
		if err == nil {
			return path
		}
	}
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return ""
	}
	return dir
}
