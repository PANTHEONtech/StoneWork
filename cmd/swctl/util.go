package main

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	installDir = filepath.Join(os.Getenv("HOME"), ".cache", "stonework", "bin")
)

type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadUrl string `json:"browser_download_url"`
	} `json:"assets"`
}

func getVppProbe() (string, error) {
	// Set the repository owner and name
	owner := "ligato"
	name := "vpp-probe"

	// Construct the URL for the latest release of the repository
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, name)

	logrus.Debugf("retrieving release: %v", url)

	// Make a GET request to the release URL to get the latest release data
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Decode the JSON response into a GitHubRelease struct
	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}

	nameOsArch := fmt.Sprintf("%s_%s", runtime.GOOS, runtime.GOARCH)

	// Find the right binary asset in the release assets
	var assetUrl string
	for _, asset := range release.Assets {
		if strings.Contains(asset.Name, nameOsArch) && strings.HasSuffix(asset.Name, ".tar.gz") {
			assetUrl = asset.BrowserDownloadUrl
			break
		}
	}

	// Return an error if no matching asset was found
	if assetUrl == "" {
		return "", fmt.Errorf("no binary asset containing %q found in release %s", nameOsArch, release.TagName)
	}

	// Construct the path to the installation file
	installPath := filepath.Join(installDir, name)
	versionPath := installPath + ".version"

	if _, err := os.Stat(installPath); err == nil {
		version, err := os.ReadFile(versionPath)
		if err == nil {
			if string(version) == release.TagName {
				logrus.Debugf("vpp-probe %s was found locally, skipping download", version)
				return installPath, nil
			} else {
				logrus.Debugf("avaiable version (%s) differs from latest release (%s), proceed to download", version, release.TagName)
			}
		} else if os.IsNotExist(err) {
			logrus.Debugf("version file not found, proceed to download")
		} else if err != nil {
			return "", err
		}
	} else if os.IsNotExist(err) {
		logrus.Debugf("vpp-probe not found, proceed to download")
	} else if err != nil {
		return "", err
	}

	// Create the installation directory if it doesn't exist
	if err := os.MkdirAll(installDir, 0755); err != nil {
		return "", err
	}

	logrus.Debugf("downloading release asset: %v", assetUrl)

	// Make a GET request to the asset URL to download the binary archive
	resp, err = http.Get(assetUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Create a new gzip reader for the response body
	gzipReader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return "", err
	}
	defer gzipReader.Close()

	// Create a new tar reader for the gzip reader
	tarReader := tar.NewReader(gzipReader)

	// Iterate over the files in the tar archive
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		// Extract the file if it is a regular file
		if header.Typeflag == tar.TypeReg && path.Base(header.Name) == "vpp-probe" {
			// Create the installation file
			installFile, err := os.Create(installPath)
			if err != nil {
				return "", fmt.Errorf("creating file failed: %w", err)
			}
			defer installFile.Close()

			// Copy the contents of the file from the tar archive to the installation file
			if _, err := io.Copy(installFile, tarReader); err != nil {
				return "", fmt.Errorf("copying file data failed: %w", err)
			}

			// Make the installation file executable
			if err := os.Chmod(installPath, 0755); err != nil {
				return "", fmt.Errorf("setting file mode failed: %w", err)
			}

			// Store the release version info
			if err := os.WriteFile(versionPath, []byte(release.TagName), 0755); err != nil {
				return "", fmt.Errorf("writing vesion to file failed: %w", err)
			}

			break
		}
	}

	return installPath, nil
}
