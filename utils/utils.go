package utils

import (
	"os"
	"path/filepath"
	"runtime"
)

// Windows save to current working directory
// other save to default temp directory
func DefaultSavePath() string {
	if runtime.GOOS == "windows" {
		ex, _ := os.Executable()
		return filepath.Dir(ex)
	} else {
		return os.TempDir()
	}
}
