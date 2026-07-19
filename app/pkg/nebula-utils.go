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

	yaml "github.com/goccy/go-yaml"
)

//TODO:Some things:
//	- the use of loos funcs + constants will get confusing but ok for now:
//	    - wrap below in receiver funcs
//	    - tidy the constants
//	- better error handling for when exec commands fail
//  - we probably don't want all of this to be in pkg but service specific backend, this will be refactored after DNS is up
//  - The UI is a mess but its fine for now as we don't have requirements yet
//  - Platform dirs should be impl
//  - Sudo outside installNebula()

type NebulaConfig struct {
	Firewall struct {
		Conntrack struct {
			DefaultTimeout string `yaml:"default_timeout"`
			TcpTimeout string `yaml:"tcp_timeout"`
			UdpTimeout string `yaml:"udp_timeout"`
		} `yaml:"conntrack"`
		Inbound []struct {
			Host string `yaml:"host"`
			Port string `yaml:"port"`
			Proto string `yaml:"proto"`
		} `yaml:"inbound"`
		InboundAction string `yaml:"inbound_action"`
		Outbound []struct {
			Host string `yaml:"host"`
			Port string `yaml:"port"`
			Proto string `yaml:"proto"`
		} `yaml:"outbound"`
		OutboundAction string `yaml:"outbound_action"`
	} `yaml:"firewall"`
	Lighthouse struct {
		AmLighthouse bool `yaml:"am_lighthouse"`
		Hosts []string `yaml:"hosts"`
		Interval int `yaml:"interval"`
	} `yaml:"lighthouse"`
	Listen struct {
		Host string `yaml:"host"`
		Port int `yaml:"port"`
	} `yaml:"listen"`
	Logging struct {
		Format string `yaml:"format"`
		Level string `yaml:"level"`
	} `yaml:"logging"`
	Pki struct {
		CaPath string `yaml:"ca"`
		CertPath string `yaml:"cert"`
		KeyPath string `yaml:"key"`
	} `yaml:"pki"`
	Punchy struct {
		Punch bool `yaml:"punch"`
	} `yaml:"punchy"`
	Relay struct {
		AmRelay bool `yaml:"am_relay"`
		UseRelays bool `yaml:"use_relays"`
	} `yaml:"relay"`
	StaticHostMap interface{} `yaml:"static_host_map"`
	Tun struct {
		Dev string `yaml:"dev"`
		Disabled bool `yaml:"disabled"`
		DropLocalBroadcast bool `yaml:"drop_local_broadcast"`
		DropMulticast bool `yaml:"drop_multicast"`
		Mtu int `yaml:"mtu"`
		Routes interface{} `yaml:"routes"`
		TxQueue int `yaml:"tx_queue"`
		UnsafeRoutes interface{} `yaml:"unsafe_routes"`
	} `yaml:"tun"`
}

func (nc *NebulaConfig) trimPath() {
	
	nc.Pki.CaPath = nc.Pki.CaPath[strings.LastIndex(nc.Pki.CaPath, "/") +1:]
	nc.Pki.CertPath = nc.Pki.CertPath[strings.LastIndex(nc.Pki.CertPath, "/") +1:]
	nc.Pki.KeyPath = nc.Pki.KeyPath[strings.LastIndex(nc.Pki.KeyPath, "/") +1:]
}

func (nc *NebulaConfig) prependHomePath() {
	nc.Pki.CaPath = filepath.Join(DESTINATION_CERTS, nc.Pki.CaPath)
	nc.Pki.CertPath = filepath.Join(DESTINATION_CERTS, nc.Pki.CertPath)
	nc.Pki.KeyPath = filepath.Join(DESTINATION_CERTS, nc.Pki.KeyPath)
}



//TODO: Add a verifier to check if path has already been fixed
//ALSO, instead of mutilating the bundle from cep server, we can just modify here until the alst slash
func InjectNebulaPath(path string) error {
	var yamlPayload NebulaConfig
	file, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("Failed to open file %s", err)
	}
	
	err = yaml.Unmarshal(file, &yamlPayload)
	
	yamlPayload.trimPath()
	yamlPayload.prependHomePath()
	
	if err != nil {
		return fmt.Errorf("Failed to extract yaml from file %s", err)
	}
	
	injectedPath, err := yaml.Marshal(yamlPayload)

	if err != nil {
		return fmt.Errorf("Failed to marshal back: %s", err)
	}

	
	err = os.WriteFile(path,  injectedPath, OWNER_READ_WRITE)

	if err != nil {
		return fmt.Errorf("Failed to write back yaml %s", err)
	}
	return nil

}

func IfNebulaExists() bool {
	_, err := os.Stat(NEBULA_PATH)
	if err != nil {
		log.Printf("nebula binary not found at %s: %s", NEBULA_PATH, err)
		return false
	}
	log.Printf("nebula binary found at %s", NEBULA_PATH)
	return true
}

func NebulaStart(nebulaPath string, certsPath string, sudoPassword string) (io.ReadCloser, error) {
	
	err := InjectNebulaPath(HOST_CONFIG)

	if err != nil {

		return nil, fmt.Errorf("failed to modify nebula config: %w", err)
	}

	cmd := exec.Command(
		"sudo", "-S", nebulaPath, "-config", filepath.Join(certsPath, "config.yml"),
	)
	cmd.Dir = certsPath

	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get stdin pipe: %w", err)
	}

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

	go func() {
		defer stdinPipe.Close()
		fmt.Fprintln(stdinPipe, sudoPassword)
	}()

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
			mode := f.Mode()
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
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
