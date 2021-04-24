package scandocs

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// LocalReadmeProvider scans filesystem for readme file
type LocalReadmeProvider struct{}

func (c *LocalReadmeProvider) GetReadme(path string) string {
	var readme []byte
	filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if !strings.EqualFold(strings.TrimSpace(f.Name()), "README.md") {
			return nil
		}

		readmef, errf := os.Open(path)
		if errf != nil {
			return nil
		}
		defer func() { readmef.Close() }()

		readme, _ = ioutil.ReadAll(readmef)
		return errors.New("found file, stopping walkign dir")
	})
	return string(readme)
}

// ReadmeScanner extracts information from readme reader
type ReadmeScanner struct{}

var deprecatedTokens = []string{
	"DEPRECATED",
	"UNMAINTAINED",
	"NOT MAINTAINED",
	"NO LONGER MAINTAINED",
	"UNSUPPORTED",
	"NOT SUPPORTED",
	"NO LONGER SUPPORTED",
}

// IsDeprecated check for signs that contents of readme says that it is deprecated
func (c *ReadmeScanner) IsDeprecated(readme string) bool {
	fs := []func(a string) string{
		func(a string) string { return a }, // identity
		strings.ToLower,
		strings.ToUpper,
	}
	for _, w := range deprecatedTokens {
		for _, f := range fs {
			if strings.Contains(readme, f(w)) {
				return true
			}
		}
	}
	return false
}
