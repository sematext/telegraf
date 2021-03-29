package processors

import (
	"os"
)

const (
	RootEnvVar     = "SPM_ROOT"
	DefaultRootDir = "/opt/spm/"
)

// GetRootDir returns the root dir of Sematext Agent installation, if it is present. Otherwise empty string.
func GetRootDir() string {
	if dir := os.Getenv(RootEnvVar); dir != "" {
		if exists(dir) {
			return dir
		}
	}

	if exists(DefaultRootDir) {
		return DefaultRootDir
	}

	return ""
}

func exists(dir string) bool {
	if _, err := os.Stat(dir); err != nil {
		return false
	}
	return true
}
