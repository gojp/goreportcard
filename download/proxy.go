package download

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
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
func ProxyDownload(path string) error {
	resp, err := http.Get(fmt.Sprintf(proxyLatestURL, path))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	var mv moduleVersion

	err = json.NewDecoder(resp.Body).Decode(&mv)
	if err != nil {
		return err
	}

	fmt.Println(mv)

	resp, err = http.Get(fmt.Sprintf(proxyZipURL, path, mv.Version))
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("status %d", resp.StatusCode)
	}

	zipPath := filepath.Base(path) + "@" + mv.Version + ".zip"
	out, err := os.Create(filepath.Join(reposDir, zipPath))
	if err != nil {
		return err
	}

	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
