package pkg

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)



type EnrolmentPayload struct {
	Config string
	CertificateAuthority string
	HostCert string
	HostKey string
}


//this only supports linux for now, for v 2.0.0 we
//would probably need to look at how fyne would work with cross-comp and save os-level operations
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

func NebulaStart() error {

	cmd := exec.Command(
	 NEBULA_PATH, "-config", filepath.Join(DESTINATION_FOLDER, "config.yaml"),
	)
	
	stderrPipe, err := cmd.StderrPipe()

	if err != nil {
		return fmt.Errorf("failed to get stderr pipe: %w", err)
	}
	
	stdoutPipe, err := cmd.StdoutPipe()
	
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}
	
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start nebula: %w", err)
	}
	
	go func() {
		stderrData, _ := io.Copy(os.Stderr, stderrPipe)
		if stderrData > 0 {
			log.Printf("nebula stderr: %s", stderrData)
		}
	}()
	
	go func() {
		stdoutData, _ := io.Copy(os.Stdout, stdoutPipe)
		if stdoutData > 0 {
			log.Printf("nebula stdout: %s", stdoutData)
		}
	}()
	
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("nebula exited with error: %w", err)
	}

	return nil
}


//NOTE: will probably not be used
func ParseCertFolder(dirPath string) (*EnrolmentPayload, error) {
	payload := &EnrolmentPayload{}
	entries, err := os.ReadDir(dirPath)
	log.Print("enrieles:, dirs", dirPath, entries)	
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}
	
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		
		name := entry.Name()
		
		switch {
		
		case name == "ca.crt":
			data, err := os.ReadFile(filepath.Join(dirPath, name))
			if err != nil {
				return nil, fmt.Errorf("failed to read ca.crt: %w", err)
			}
			payload.CertificateAuthority = string(data)
		
		case strings.HasSuffix(name, ".key"):
			data, err := os.ReadFile(filepath.Join(dirPath, name))
			if err != nil {
				return nil, fmt.Errorf("failed to read %s: %w", name, err)
			}
			payload.HostKey = string(data)
		
		case strings.HasSuffix(name, ".crt") && name != "ca.crt":
			data, err := os.ReadFile(filepath.Join(dirPath, name))
			if err != nil {
				return nil, fmt.Errorf("failed to read %s: %w", name, err)
			}
			payload.HostCert = string(data)
		
		case name == "config.yaml":
			data, err := os.ReadFile(filepath.Join(dirPath, name))
			if err != nil {
				return nil, fmt.Errorf("failed to read config.yaml: %w", err)
			}
			payload.Config = string(data)
		}
	}
	return payload, nil
}

//Copy pasted
//TODO: dump contents of each unzipped file into path obj
//copme back to this
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
        if !strings.HasPrefix(path, filepath.Clean(dest) + string(os.PathSeparator)) {
            return fmt.Errorf("illegal file path: %s", path)
        }

        if f.FileInfo().IsDir() {
            os.MkdirAll(path, f.Mode())
        } else {
            os.MkdirAll(filepath.Dir(path), f.Mode())
            f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
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

func GetHomePath() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		log.Printf("error when fetching home path: %s", err)
		return ""
	}
	return dir
}
