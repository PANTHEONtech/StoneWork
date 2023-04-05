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
	"time"

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

func downloadVppProbe() (string, error) {
	const (
		repoOwner = "ligato"
		repoName  = "vpp-probe"

		latestVersionCheckPeriod = time.Hour
	)

	// Construct the path to the installation file
	installPath := filepath.Join(installDir, repoName)
	versionPath := installPath + ".version"

	var installedVersion string
	var lastCheck time.Time

	if _, err := os.Stat(installPath); err == nil {
		info, err := os.Stat(versionPath)
		if err == nil {
			lastCheck = info.ModTime()
		}
		version, err := os.ReadFile(versionPath)
		if err == nil {
			installedVersion = string(version)
		} else if os.IsNotExist(err) {
			logrus.Debugf("version file not found, proceed to download")
		} else if err != nil {
			return "", err
		}
	} else if os.IsNotExist(err) {
		logrus.Debugf("vpp-probe install directory not found, proceed to download")
	} else if err != nil {
		return "", err
	}

	if installedVersion != "" {
		logrus.Debugf("installed version of vpp-probe: %v", installedVersion)
		if d := time.Since(lastCheck); d < latestVersionCheckPeriod {
			logrus.Debugf("last check or download occurred recently %v ago (less than %v), skipping check for the latest release", d.Round(time.Minute), latestVersionCheckPeriod)
			return installPath, nil
		}
	}

	// Construct the URL for the latest release of the repository
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", repoOwner, repoName)

	logrus.Debugf("checking latest release: %v", url)

	// Make a GET request to the release URL to get the latest release data
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve latest release from GitHub: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode > 299 {
		return "", fmt.Errorf("response failed with status code: %d (%s) and\nbody: %s\n", resp.StatusCode, resp.Status, body)
	}
	if err != nil {
		return "", fmt.Errorf("failed reading response body: %w", err)
	}

	// Decode the JSON response into a GitHubRelease struct
	var release GitHubRelease
	if err := json.Unmarshal(body, &release); err != nil {
		return "", fmt.Errorf("failed to decode data for latest release: %w", err)
	}

	latestVersion := release.TagName
	logrus.Debugf("latest release version on GitHub: %v", latestVersion)

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
		return "", fmt.Errorf("no binary asset containing %q found in release %s", nameOsArch, latestVersion)
	}

	// Compare installed version with the latest release
	if installedVersion == latestVersion {
		logrus.Debugf("vpp-probe is already at latest version %s, skipping download", installedVersion)
		if err := os.WriteFile(versionPath, []byte(latestVersion), 0755); err != nil {
			return "", fmt.Errorf("writing version to file failed: %w", err)
		}
		return installPath, nil
	} else {
		logrus.Debugf("vpp-probe version (%s) differs from latest release (%s), proceed to download", installedVersion, latestVersion)
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

			if err := extractVppProbe(tarReader, installPath); err != nil {
				return "", fmt.Errorf("failed to extract vpp-probe binary: %w", err)
			}

			// Store the release version info
			if err := os.WriteFile(versionPath, []byte(latestVersion), 0755); err != nil {
				return "", fmt.Errorf("writing version to file failed: %w", err)
			}

			break
		}
	}

	return installPath, nil
}

func extractVppProbe(tarReader io.Reader, installPath string) error {
	// Create the installation file
	installFile, err := os.Create(installPath)
	if err != nil {
		return fmt.Errorf("creating file failed: %w", err)
	}
	defer installFile.Close()

	// Copy the contents of the file from the tar archive to the installation file
	if _, err := io.Copy(installFile, tarReader); err != nil {
		return fmt.Errorf("copying file data failed: %w", err)
	}

	// Make the installation file executable
	if err := os.Chmod(installPath, 0755); err != nil {
		return fmt.Errorf("setting file mode failed: %w", err)
	}

	return nil
}
