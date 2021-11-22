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
	proxyModURL    = "https://proxy.golang.org/%s/@v/%s.mod"
	reposDir       = "_repos/src"
)

type moduleVersion struct {
	Version string
}

// ModuleName gets the name of a module from the proxy
func ModuleName(path string) (string, error) {
	lowerPath := strings.ToLower(path)

	ver, err := LatestVersion(path)
	if err != nil {
		return "", err
	}

	u := fmt.Sprintf(proxyModURL, lowerPath, ver)
	resp, err := http.Get(u)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("could not get latest module version from %s: %s", u, string(b))
	}

	sp := strings.Split(string(b), "\n")
	if len(sp) == 0 {
		return "", fmt.Errorf("empty go.mod")
	}

	mn := strings.Fields(sp[0])
	if len(mn) != 2 {
		return "", fmt.Errorf("invalid go.mod: %s", mn)
	}

	return mn[1], nil
}

// LatestVersion gets the latest module version from the proxy
func LatestVersion(path string) (string, error) {
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

	return mv.Version, nil
}

// ProxyDownload downloads a package from proxy.golang.org
func ProxyDownload(path string) (string, error) {
	lowerPath := strings.ToLower(path)

	ver, err := LatestVersion(path)
	if err != nil {
		return "", err
	}

	resp, err := http.Get(fmt.Sprintf(proxyZipURL, lowerPath, ver))
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("status %d", resp.StatusCode)
	}

	zipPath := filepath.Join(reposDir, filepath.Base(path)+"@"+ver+".zip")
	out, err := os.Create(zipPath)
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

	cmd := exec.Command("unzip", "-o", zipPath, "-d", reposDir)

	err = cmd.Run()
	if err != nil {
		return "", err
	}

	err = os.RemoveAll(zipPath)
	if err != nil {
		return "", err
	}

	err = os.Rename(filepath.Join(reposDir, lowerPath+"@"+ver), filepath.Join(reposDir, lowerPath))
	if err != nil {
		return "", err
	}

	return ver, nil
}
