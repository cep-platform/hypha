package pkg

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const nebulaReleaseBaseURL = "https://github.com/slackhq/nebula/releases/download"

// buildDownloadURL returns the release archive URL and whether it is a tar.gz (vs zip).
// Nebula release naming conventions:
//   - Linux:   nebula-linux-{arch}.tar.gz  (amd64, arm64, …)
//   - Darwin:  nebula-darwin.zip           (universal fat binary)
//   - Windows: nebula-windows-{arch}.zip   (amd64, arm64, …)
func buildDownloadURL(version string) (url string, isTarGz bool, err error) {
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	var filename string
	switch goos {
	case "linux":
		filename = fmt.Sprintf("nebula-linux-%s.tar.gz", goarch)
		isTarGz = true
	case "darwin":
		filename = "nebula-darwin.zip"
	case "windows":
		filename = fmt.Sprintf("nebula-windows-%s.zip", goarch)
	default:
		return "", false, fmt.Errorf("unsupported OS: %s", goos)
	}

	url = fmt.Sprintf("%s/v%s/%s", nebulaReleaseBaseURL, version, filename)
	return url, isTarGz, nil
}

// InstallNebula downloads the Nebula binary for the current OS/arch from the
// official GitHub releases page and extracts it to the directory defined by NEBULA_PATH.
func InstallNebula() error {
	destDir := filepath.Dir(NEBULA_PATH)

	url, isTarGz, err := buildDownloadURL(NEBULA_VERSION)
	if err != nil {
		return fmt.Errorf("could not determine download URL: %w", err)
	}

	log.Printf("downloading nebula v%s from %s", NEBULA_VERSION, url)

	tmpFile, err := downloadToTemp(url)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	defer os.Remove(tmpFile)

	if err := os.MkdirAll(destDir, DEFAULT_PERMISSIONS); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	if isTarGz {
		if err := Untar(tmpFile, destDir); err != nil {
			return fmt.Errorf("extraction failed: %w", err)
		}
	} else {
		if err := Unzip(tmpFile, destDir); err != nil {
			return fmt.Errorf("extraction failed: %w", err)
		}
	}

	log.Printf("nebula v%s installed to %s", NEBULA_VERSION, destDir)
	return nil
}

// downloadToTemp streams the content at url into a temporary file and returns its path.
func downloadToTemp(url string) (string, error) {
	resp, err := http.Get(url) //nolint:gosec // URL is constructed internally from a trusted constant base
	if err != nil {
		return "", fmt.Errorf("HTTP GET failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected HTTP status %d for %s", resp.StatusCode, url)
	}

	tmp, err := os.CreateTemp("", "nebula-download-*")
	if err != nil {
		return "", fmt.Errorf("could not create temp file: %w", err)
	}
	defer tmp.Close()

	if _, err := io.Copy(tmp, resp.Body); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("failed to write download to disk: %w", err)
	}

	return tmp.Name(), nil
}

// Untar extracts a .tar.gz archive into dest, stripping any leading directory component
// so the binaries land directly in dest.
func Untar(src, dest string) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return fmt.Errorf("not a valid gzip file: %w", err)
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	cleanDest := filepath.Clean(dest) + string(os.PathSeparator)

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("tar read error: %w", err)
		}

		// Strip leading directory component (e.g. "nebula-v1.9.5-linux-amd64/nebula" → "nebula")
		parts := strings.SplitN(hdr.Name, "/", 2)
		name := hdr.Name
		if len(parts) == 2 {
			name = parts[1]
		}
		if name == "" {
			continue
		}

		target := filepath.Join(dest, name)

		// Guard against zip-slip / path traversal
		if !strings.HasPrefix(filepath.Clean(target)+string(os.PathSeparator), cleanDest) &&
			filepath.Clean(target) != filepath.Clean(dest) {
			return fmt.Errorf("illegal path in archive: %s", hdr.Name)
		}

		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, DEFAULT_PERMISSIONS); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), DEFAULT_PERMISSIONS); err != nil {
				return err
			}
			out, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(hdr.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(out, tr); err != nil {
				out.Close()
				return err
			}
			out.Close()
		}
	}

	return nil
}
