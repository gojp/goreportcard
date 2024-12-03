package download

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	reposDir = "_repos/src"
)

type moduleVersion struct {
	Version string
}

// ProxyClient is a client for the module proxy
type ProxyClient struct {
	URL string
}

// NewProxyClient returns a new ProxyClient
func NewProxyClient(url string) ProxyClient {
	return ProxyClient{URL: url}
}

func (c *ProxyClient) latestURL(module string) string {
	return fmt.Sprintf("%s/%s/@latest", c.URL, module)
}

func (c *ProxyClient) zipURL(module, version string) string {
	return fmt.Sprintf("%s/%s/@v/%s.zip", c.URL, module, version)
}

func (c *ProxyClient) modURL(module, version string) string {
	return fmt.Sprintf("%s/%s/@v/%s.mod", c.URL, module, version)
}

// ModuleName gets the name of a module from the proxy
func (c *ProxyClient) ModuleName(path string) (string, error) {
	lowerPath := strings.ToLower(path)

	ver, err := c.LatestVersion(path)
	if err != nil {
		return "", err
	}

	u := c.modURL(lowerPath, ver)
	resp, err := http.Get(u)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("could not get module name from %s", u)
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
func (c *ProxyClient) LatestVersion(path string) (string, error) {
	lowerPath := strings.ToLower(path)
	u := c.latestURL(lowerPath)
	resp, err := http.Get(u)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("could not get latest module version from %s", u)
	}

	var mv moduleVersion

	err = json.NewDecoder(resp.Body).Decode(&mv)
	if err != nil {
		return "", err
	}

	return mv.Version, nil
}

// ProxyDownload downloads a package from proxy.golang.org
func (c *ProxyClient) ProxyDownload(path string) (string, error) {
	lowerPath := strings.ToLower(path)

	ver, err := c.LatestVersion(path)
	if err != nil {
		return "", err
	}

	resp, err := http.Get(c.zipURL(lowerPath, ver))
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

	// err = os.RemoveAll(filepath.Join(reposDir, path, "@"+ver))
	// if err != nil {
	// 	return "", err
	// }

	cmd := exec.Command("unzip", "-o", zipPath, "-d", reposDir)

	err = cmd.Run()
	if err != nil {
		return "", err
	}

	err = os.RemoveAll(zipPath)
	if err != nil {
		return "", err
	}

	// err = os.Rename(filepath.Join(reposDir, lowerPath+"@"+ver), filepath.Join(reposDir, lowerPath))
	// if err != nil {
	// 	return "", err
	// }

	return ver, nil
}
