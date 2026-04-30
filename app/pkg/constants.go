package pkg

import (
	"path/filepath"	
)
var (
	HOME_DIR = GetHomeDir()
	
	NEBULA_PATH = filepath.Join(
		HOME_DIR, ".cache", "nebula", "nebula",
		)


	ZIPPED_CERT_PATH = filepath.Join(
		 HOME_DIR, ".config", "nebula-certs",
		)
		
	DESTINATION_FOLDER = filepath.Join(
	HOME_DIR, "/.cache",
	)
	HOST_NAME = "nebula"	
	HOST_PATH = filepath.Join(ZIPPED_CERT_PATH, HOST_NAME + ".cepbundle")


	DESTINATION_CERTS = filepath.Join(
		DESTINATION_FOLDER, HOST_NAME,
		)
	
	DIRS = []string{
		ZIPPED_CERT_PATH,
		DESTINATION_FOLDER,
	}
)

const (
	DEFAULT_PERMISSIONS = 0755
	OWNER_READ_WRITE    = 0644
	NEBULA_VERSION      = "1.9.5"
)

