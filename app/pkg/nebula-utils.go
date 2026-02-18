package pkg

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

//TODO:Some things:
//	- the use of loos funcs + constants will get confusing but ok for now:
//	    - wrap below in receiver funcs
//	    - tidy the constants
//	- better error handling for when exec commands fail
//  - zip filepath needs to be given by user and validated 
//  - we probably don't want all of this to be in pkg but service specific backend, this will be refactored after DNS is up
//  - The UI is a mess but its fine for now as we don't have requirements yet

// this only supports linux for now, for v 2.0.0 we
// would probably need to look at how fyne would work with cross-comp and save os-level operations
func InstallNebula() error {
	cmd := exec.Command(
		"sudo", "pacman", "install", "nebula",
	)

	//buffers for debugging
	var stdoutBuff bytes.Buffer
	var stderrBuff bytes.Buffer

	cmd.Stdout = &stdoutBuff
	cmd.Stderr = &stderrBuff

	err := cmd.Run()

	if err != nil {
		stderrString := stderrBuff.String()
		log.Printf("error encountered when installing nebula: %s", stderrString)
		return fmt.Errorf("error encountered when installing nebula: %s ", stderrString)
	}

	log.Printf("Nebula successfully installed: %s", stdoutBuff.String())
	return nil
}

func IfNebulaExists() bool {
	cmd := exec.Command("nebula", "-version")

	//buffers for debugging
	var stdoutBuff bytes.Buffer
	var stderrBuff bytes.Buffer

	cmd.Stdout = &stdoutBuff
	cmd.Stderr = &stderrBuff

	err := cmd.Run()

	if err != nil {
		stderrString := stderrBuff.String()
		log.Printf("nebula not installed: %s", stderrString)
		return false
	}

	log.Printf("nebula successfully installed: %s", stdoutBuff.String())
	return true
}

func NebulaStart(nebulaPath string, certsPath string) (io.ReadCloser, error) {

	cmd := exec.Command(
		"sudo", nebulaPath, "-config", filepath.Join(certsPath, "config.yaml"),
	)

	stderrPipe, err := cmd.StderrPipe()

	if err != nil {
		return stderrPipe, fmt.Errorf("failed to get stderr pipe: %w", err)
	}

	stdoutPipe, err := cmd.StdoutPipe()

	if err != nil {
		return stderrPipe, fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return stderrPipe, fmt.Errorf("failed to start nebula: %w", err)
	}

	go func() {
		stderrData, _ := io.Copy(os.Stderr, stderrPipe)
		if stderrData > 0 {
			log.Printf("nebula stderr: %v", stderrData)
		}
	}()

	go func() {
		stdoutData, _ := io.Copy(os.Stdout, stdoutPipe)
		if stdoutData > 0 {
			log.Printf("nebula stdout: %v", stdoutData)
		}
	}()

	return stdoutPipe, nil
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
