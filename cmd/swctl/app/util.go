package app

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"path"
	"runtime"
	"strings"
)

var defaultVppProbeEnv = "docker"

type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadUrl string `json:"browser_download_url"`
	} `json:"assets"`
}

func downloadAndExtractSubAsset(assetUrl string, subAssetName string, extractionFilePath string) error {
	logrus.Debugf("downloading release asset: %v", assetUrl)

	// Make a GET request to the asset URL to download the binary archive
	resp, err := http.Get(assetUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create a new gzip reader for the response body
	gzipReader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return err
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
			return err
		}

		// Extract the file if it is a regular file
		if header.Typeflag == tar.TypeReg && path.Base(header.Name) == subAssetName {
			if err := extractFromReader(tarReader, extractionFilePath); err != nil {
				return fmt.Errorf("failed to extract %s binary: %w", subAssetName, err)
			}
			return nil
		}
	}

	return fmt.Errorf("couldn't find the subasset %s in downloaded asset (url: %s)", subAssetName, assetUrl)
}

func retrieveReleaseAssetUrl(repoOwner string, repoName string, releaseCommitTag string) (string, error) {
	// Construct the URL for the given release of the repository
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/tags/%s", repoOwner, repoName, releaseCommitTag)
	logrus.Debugf("checking %s release: %v", releaseCommitTag, url)

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

	releaseVersion := release.TagName
	logrus.Debugf("target release version on GitHub: %v", releaseVersion)

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
		return "", fmt.Errorf("no binary asset containing %q found in release %s", nameOsArch, releaseVersion)
	}
	return assetUrl, nil
}

func extractFromReader(reader io.Reader, extractionFilePath string) error {
	// Create the installation file
	installFile, err := os.Create(extractionFilePath)
	if err != nil {
		return fmt.Errorf("creating file failed: %w", err)
	}
	defer installFile.Close()

	// Copy the contents of the file from the tar archive to the installation file
	if _, err := io.Copy(installFile, reader); err != nil {
		return fmt.Errorf("copying file data failed: %w", err)
	}

	// Make the installation file executable
	if err := os.Chmod(extractionFilePath, 0755); err != nil {
		return fmt.Errorf("setting file mode failed: %w", err)
	}

	return nil
}
