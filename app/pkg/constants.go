package pkg

import (
	"path/filepath"	
)
var (
	HOME_DIR = GetHomePath()
	
	NEBULA_PATH = filepath.Join(
		HOME_DIR, ".cache", "nebula", "nebula",
		)


	ZIPPED_CERT_PATH = filepath.Join(
		HOME_DIR, ".config", "nebula-certs",
		)
	DESTINATION_FOLDER = filepath.Join(
		HOME_DIR, "etc", "nebula",
	)

	DESTINATION_CERTS = filepath.Join(
		DESTINATION_FOLDER, "host",
		)
	
	DIRS = []string{
		ZIPPED_CERT_PATH,
		DESTINATION_FOLDER,
	}
	
	HOST_NAME = "nebula"
	HOST_PATH = filepath.Join(ZIPPED_CERT_PATH, HOST_NAME + ".zip")
	
)

const (
	DEFAULT_PERMISSIONS  = 0755
)

