package gitstats

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path"
)

// GitCmdLocalClient works with local git through os commands
type GitCmdLocalClient struct {
	Path string
}

// Clone git repo
func (g *GitCmdLocalClient) Clone(gitURL url.URL) error {
	dirPath := g.DirPath(gitURL)
	if _, err := os.Stat(dirPath); !os.IsNotExist(err) {
		return nil
	}
	return exec.Command("git", "clone", gitURL.String(), dirPath).Run()
}

// GetGitLog fetches git log entries given path for git
func (g *GitCmdLocalClient) GetGitLog(gitURL url.URL) (GitLog, error) {
	cmd := exec.Command(
		"git",
		fmt.Sprintf("--git-dir=%s/.git", g.DirPath(gitURL)),
		"log",
		"--pretty=format:%at %ae",
	)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("can not get stdout pipe: %w", err)
	}
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("can not start go command: %w", err)
	}
	scanner := bufio.NewScanner(stdout)

	var gitlogs []GitLogEntry
	for scanner.Scan() {
		entry, err := NewGitLogEntryFromLine(scanner.Text())
		if err != nil {
			return nil, fmt.Errorf("issue parsing git log line: %w", err)
		}
		gitlogs = append(gitlogs, entry)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("got error from stdout of git log scanner: %w", err)
	}
	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("command did not finish successfully: %w", err)
	}
	return gitlogs, nil
}

// DirPath gets path where git repo is stored locally
func (g *GitCmdLocalClient) DirPath(gitURL url.URL) string {
	return path.Join(g.Path, dirName(gitURL))
}

// dirName returns safe name for git URL.
// This is easy way to avoid special characters that can be in URL.
func dirName(repoURL url.URL) string {
	return base64.URLEncoding.EncodeToString([]byte(repoURL.String()))
}
