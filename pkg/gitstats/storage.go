package gitstats

import (
	"crypto/md5"
	"encoding/hex"
	"net/url"
	"os"
	"os/exec"
	"path"
)

// GitProcessStorage provides utilities to fetch git repositories utilizing system git command
type GitProcessStorage struct {
	Path string
}

// Fetch invokes git clone system command to folder
func (g *GitProcessStorage) Fetch(gitURL url.URL) error {
	if _, err := os.Stat(g.DirPath(gitURL)); !os.IsNotExist(err) {
		return nil
	}
	return exec.Command("git", "clone", gitURL.String(), g.DirPath(gitURL)).Run()
}

// DirPath gets path where git repo is stored locally
func (g *GitProcessStorage) DirPath(gitURL url.URL) string {
	return path.Join(g.Path, dirName(gitURL))
}

func dirName(repoURL url.URL) string {
	v := md5.Sum([]byte(repoURL.String()))
	return hex.EncodeToString(v[:])
}
