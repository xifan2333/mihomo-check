package utils

import (
	"os"
	"path/filepath"
)

func GetExecutablePath() string {
	ex, err := os.Executable()
	if err != nil {
		LogError("get executable path failed: %v", err)
		return "."
	}
	return filepath.Dir(ex)
}
