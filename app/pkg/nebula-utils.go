package pkg

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"sync"
)

//TODO:Some things:
//	- the use of loos funcs + constants will get confusing but ok for now:
//	    - wrap below in receiver funcs
//	    - tidy the constants
//	- better error handling for when exec commands fail
//  - zip filepath needs to be given by user and validated 
//  - we probably don't want all of this to be in pkg but service specific backend, this will be refactored after DNS is up
//  - The UI is a mess but its fine for now as we don't have requirements yet
//  - Platform dirs should be impl
//  - Sudo outside installNebula()


func IfNebulaExists() bool {
	_, err := os.Stat(NEBULA_PATH)
	if err != nil {
		log.Printf("nebula binary not found at %s: %s", NEBULA_PATH, err)
		return false
	}
	log.Printf("nebula binary found at %s", NEBULA_PATH)
	return true
}

func NebulaStart(nebulaPath string, certsPath string, password string) (io.ReadCloser, error) {
	// -S tells sudo to read the password from stdin instead of a terminal
	cmd := exec.Command(
		"sudo", nebulaPath, "-config", filepath.Join(certsPath, "config.yml"),
	)
	cmd.Stdin = strings.NewReader(password + "\n")

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start nebula: %w", err)
	}

	// Fan stdout and stderr into a single reader so the UI sees all output.
	pr, pw := io.Pipe()
	var wg sync.WaitGroup
	wg.Add(2)

	copy := func(src io.Reader) {
		defer wg.Done()
		io.Copy(pw, src) //nolint:errcheck
	}

	go copy(stdoutPipe)
	go copy(stderrPipe)

	go func() {
		wg.Wait()
		pw.Close()
	}()

	return pr, nil
}

func Unzip(src string, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, DEFAULT_PERMISSIONS)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()
		path := filepath.Join(dest, f.Name)
		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, DEFAULT_PERMISSIONS)
		} else {
			os.MkdirAll(filepath.Dir(path), DEFAULT_PERMISSIONS)
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, OWNER_READ_WRITE)
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		// add to payload obj here

		return nil
	}
	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}

func ValidateDir(dirs []string) error {
	for _, dir := range dirs {
		err := os.MkdirAll(dir, DEFAULT_PERMISSIONS)
		if err != nil && !os.IsExist(err) {
			log.Printf("error when trying to make %s directory: %w", dir, err.Error())
			return err
		}
	}
	return nil
}

func GetHomeDir() string {
	sudoUser := os.Getenv("SUDO_USER")
	if sudoUser != "" {
		// Running with sudo, get the real user's home
		usr, err := user.Lookup(sudoUser)
		if err == nil {
			return usr.HomeDir
		}
	}
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	return usr.HomeDir
}
