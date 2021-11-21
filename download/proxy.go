package download

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	proxyLatestURL = "https://proxy.golang.org/%s/@latest"
	proxyZipURL    = "https://proxy.golang.org/%s/@v/%s.zip"
	reposDir       = "_repos/src"
)

type moduleVersion struct {
	Version string
}

// ProxyDownload downloads a package from proxy.golang.org
func ProxyDownload(path string) (string, error) {
	lowerPath := strings.ToLower(path)
	u := fmt.Sprintf(proxyLatestURL, lowerPath)
	resp, err := http.Get(u)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("could not get latest module version from %s: %s", u, string(b))
	}

	var mv moduleVersion

	err = json.NewDecoder(resp.Body).Decode(&mv)
	if err != nil {
		return "", err
	}

	resp, err = http.Get(fmt.Sprintf(proxyZipURL, lowerPath, mv.Version))
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("status %d", resp.StatusCode)
	}

	zipPath := filepath.Base(path) + "@" + mv.Version + ".zip"
	out, err := os.Create(filepath.Join(reposDir, zipPath))
	if err != nil {
		return "", err
	}

	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	err = os.RemoveAll(filepath.Join(reposDir, path))
	if err != nil {
		return "", err
	}

	cmd := exec.Command("unzip", "-o", filepath.Join(reposDir, zipPath), "-d", reposDir)

	err = cmd.Run()
	if err != nil {
		return "", err
	}

	err = os.RemoveAll(filepath.Join(reposDir, zipPath))
	if err != nil {
		return "", err
	}

	err = os.Rename(strings.ToLower(filepath.Join(reposDir, path+"@"+mv.Version)), strings.ToLower(filepath.Join(reposDir, path)))
	if err != nil {
		return "", err
	}

	return mv.Version, nil
}
