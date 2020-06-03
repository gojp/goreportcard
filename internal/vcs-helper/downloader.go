package vcshelper

type IDownloader interface {
	// Download to copy git repo to local folders
	// @repoRoot is relative path to repo
	// @err error
	Download(remoteURL string, localDir string) (repoRoot string, err error)
}
