package vcshelper

import (
	"path/filepath"

	"github.com/gojp/goreportcard/internal/helper"
	"github.com/pkg/errors"
	"github.com/yeqown/log"
	"golang.org/x/tools/go/vcs"
)

var _ IDownloader = builtinToolVCS{}

type builtinToolVCS struct {
	*vcs.Cmd

	gitPrefixes map[string]string // map[host]prefix
}

// NewBuiltinToolVCS .
func NewBuiltinToolVCS(cfgs map[string]string) IDownloader {
	downloader := builtinToolVCS{
		Cmd:         vcs.ByCmd("git"),
		gitPrefixes: make(map[string]string, 4),
	}

	for k, v := range cfgs {
		downloader.gitPrefixes[k] = v
	}

	return downloader
}

// Download .
func (btvcs builtinToolVCS) Download(repoURL string, localDir string) (repoRoot string, err error) {
	outs, err := hdlRepoURL(repoURL)
	if err != nil {
		log.Errorf("could hdl repoURL=%s, err=%v", repoURL, err)
		return localDir, errors.Wrap(err, "gitDownload.clone failed to hdlRepoURL")
	}
	host, owner, repoName := outs[0], outs[1], outs[2]

	// make sure the path has exists
	localDir = filepath.Join(localDir, host, owner, repoName)
	if err := helper.EnsurePath(localDir); err != nil {
		return localDir, errors.Wrap(err, "gitDownload.clone.EnsurePath")
	}

	// get git prefix
	var (
		prefix string
		ok     bool
	)
	if prefix, ok = btvcs.gitPrefixes[host]; !ok {
		return "", errors.New("gitDownload.clone no such host config")
	}
	remoteURI := wrapRepoURL(prefix, host, owner, repoName)

	// FIXME: check path exist or not, if exists then using download
	if err := btvcs.Cmd.Create(localDir, remoteURI); err != nil {
		log.Warnf("builtinToolVCS.Download failed to Cmd.Create, err=%v", err)
		return localDir, err
	}

	return localDir, nil
}
