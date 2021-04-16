package gitstats

import (
	"bufio"
	"fmt"
	"math"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// GitLogEntry contains info about single git log entry
type GitLogEntry struct {
	AuthorEmail string
	AuthorDate  time.Time
}

func (e GitLogEntry) MonthsSinceLastCommit() uint {
	return uint(math.Floor(float64(time.Since(e.AuthorDate).Hours()) / 24 / 28))
}

// NewGitLogEntryFromLine git log entry from single row of predefined pretty print format
func NewGitLogEntryFromLine(input string) (GitLogEntry, error) {
	vals := strings.Split(strings.TrimSpace(input), " ")
	if len(vals) != 2 {
		return GitLogEntry{}, fmt.Errorf("wrong number of args for string: %s", input)
	}

	createdAt, err := strconv.ParseInt(vals[0], 10, 64)
	if err != nil {
		return GitLogEntry{}, fmt.Errorf("bad UNIX timestamp format for string(%s): %w", vals[0], err)
	}

	entry := GitLogEntry{
		AuthorEmail: vals[1],
		AuthorDate:  time.Unix(createdAt, 0),
	}
	return entry, nil
}

// GetLog fetches git log entries given path for git
func GetLog(gitPath string) ([]GitLogEntry, error) {
	cmd := exec.Command(
		"git",
		fmt.Sprintf("--git-dir=%s/.git", gitPath),
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
