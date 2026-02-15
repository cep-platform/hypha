package pkg

import (
	"path/filepath"	
)
var (
	HOME_DIR = GetHomePath()
	ZIPPED_CERT_PATH = filepath.Join(
		HOME_DIR, ".config", "nebula-certs",
		)
	DESTINATION_CERT_PATH = filepath.Join(
		HOME_DIR, ".cache", "unzipped-certificates",
		)
	DIRS = []string{
		ZIPPED_CERT_PATH,
		DESTINATION_CERT_PATH,
	}
	HOST_NAME = "host"
	HOST_PATH = filepath.Join(ZIPPED_CERT_PATH, HOST_NAME + ".zip")
)

const (
	DEFAULT_PERMISSIONS  = 0755
)

