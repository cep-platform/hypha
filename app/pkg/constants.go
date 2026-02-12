package pkg

import (
	"path/filepath"	
)
var (
	HOME_DIR = GetHomePath()
	ZIPPED_CERT_PATH = filepath.Dir(filepath.Join(
		HOME_DIR, ".config", "nebula-certs"),
	)
	DESTINATION_CERT_PATH = filepath.Dir(filepath.Join(
		HOME_DIR, ".cache", "unzipped-certificates"),
		)
)

