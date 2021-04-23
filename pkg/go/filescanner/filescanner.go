package filescanner

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type FileScanner struct{}

func (f *FileScanner) HasTests(path string) bool {
	found := false
	filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if !strings.HasSuffix(f.Name(), "_test.go") {
			return nil
		}
		if has, _ := fileHasString(path, "func Test"); has {
			found = has
			return errors.New("found benchmark")
		}
		return nil
	})
	return found
}

func (f *FileScanner) HasBenchmarks(path string) bool {
	found := false
	filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if !strings.HasSuffix(f.Name(), "_test.go") {
			return nil
		}
		if has, _ := fileHasString(path, "func Bench"); has {
			found = has
			return errors.New("found benchmark")
		}
		return nil
	})
	return found
}

func fileHasString(path string, target string) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return false, fmt.Errorf("can not open file: %w", err)
	}
	defer func() { file.Close() }()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), target) {
			return true, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return false, fmt.Errorf("got error from scanner: %w", err)
	}
	return false, nil
}
