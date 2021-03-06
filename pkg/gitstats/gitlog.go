package gitstats

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// GitLogEntry contains info about single git log entry
type GitLogEntry struct {
	AuthorEmail string
	AuthorDate  time.Time
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

// GitLog is sequence of git log entries in reverse chronological order (i.e. first is latest)
type GitLog []GitLogEntry

// NumContributors returns number of contributors in log
func (logs GitLog) NumContributors() uint {
	var count uint = 0
	contributors := map[string]bool{}
	for _, entry := range logs {
		if !contributors[entry.AuthorEmail] {
			contributors[entry.AuthorEmail] = true
			count++
		}
	}
	return count
}

// DaysSinceLastCommit returns partial days since last commit
func (logs GitLog) DaysSinceLastCommit() float64 {
	if len(logs) == 0 {
		return 0
	}
	return time.Since(logs[0].AuthorDate).Hours() / 24
}
