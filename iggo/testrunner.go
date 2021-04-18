package iggo

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

type GoCmdTestRunner struct{}

// GoModuleTestRunResult is summary of running tests in for all packages in Go module
type GoModuleTestRunResult struct {
	HasTests                  bool    `json:"has_tests"`
	HasTestFiles              bool    `json:"has_test_files"`
	NumPackages               uint32  `json:"num_packages"`
	NumPackagesWithTests      uint32  `json:"num_packages_with_tests"`
	NumPackagesWithTestsFiles uint32  `json:"num_packages_with_tests_files"`
	NumPackagesTestsPassed    uint32  `json:"nun_packages_tests_passed"`
	MinPackageCoverage        float64 `json:"min_package_coverage"`
	AvgPackageCoverage        float64 `json:"avg_package_coverage"`
}

// RunModuleTets runs tests for all packages in Go module, collects aggregate statistics
func (c *GoCmdTestRunner) RunModuleTets(moduleDirPath string) (GoModuleTestRunResult, error) {
	pkgTestStats, err := c.RunTests(moduleDirPath)
	if err != nil {
		return GoModuleTestRunResult{}, err
	}

	stats := GoModuleTestRunResult{}

	sumCov := 0.0
	for _, v := range pkgTestStats {
		if v.HasTests {
			stats.HasTests = true
			stats.NumPackagesWithTests++
		}
		if v.HasTestFiles {
			stats.HasTestFiles = true
			stats.NumPackagesWithTestsFiles++
		}

		stats.NumPackages++
		if v.Passed {
			stats.NumPackagesTestsPassed++
		}
		if v.Coverage > 0 && (stats.MinPackageCoverage == 0 || v.Coverage < stats.MinPackageCoverage) {
			stats.MinPackageCoverage = v.Coverage
		}
		sumCov += v.Coverage
	}

	stats.AvgPackageCoverage = sumCov / float64(stats.NumPackagesWithTests)
	return stats, nil
}

type GoPackageTestRunResult struct {
	Package      string  `json:"-"`
	Passed       bool    `json:"passed"`
	HasTests     bool    `json:"has_tests"`
	HasTestFiles bool    `json:"has_tests_files"`
	Coverage     float64 `json:"coverage"`
}

// RunTests runs tests via Go process and returns report
func (c *GoCmdTestRunner) RunTests(moduleDirPath string) (map[string]GoPackageTestRunResult, error) {
	cmd := exec.Command("go", "test", "-json", "-covermode=atomic", "./...")
	cmd.Dir = moduleDirPath
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("can not get stdout pipe: %w", err)
	}
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("can not start go command: %w", err)
	}
	scanner := bufio.NewScanner(stdout)

	results := map[string]GoPackageTestRunResult{}
	for scanner.Scan() {

		var row GoPackageTestRunResultOutput
		if err := json.Unmarshal(scanner.Bytes(), &row); err != nil {
			// intentionally ignoring error, since ignoring lines that do not match format
			log.Println(err)
			continue
		}

		// same package can be in multiple rows with different info
		// check it is not added yet and add
		if _, ok := results[row.Package]; !ok {
			results[row.Package] = GoPackageTestRunResult{Package: row.Package}
		}

		// get old version of package result, not can not mutate struct in map in Go
		v := results[row.Package]

		switch {
		case row.Action == GoTestResultActionPass:
			v.Passed = true
			v.HasTestFiles = true
			v.HasTests = true
		case row.Action == GoTestResultActionFail:
			v.Passed = false
			v.HasTestFiles = true
			v.HasTests = true
		case row.Action == GoTestResultActionOutput && strings.Contains(row.Output, "testing: warning: no tests to run"):
			v.HasTestFiles = true
			v.HasTests = false
		case row.Action == GoTestResultActionOutput && strings.HasPrefix(row.Output, "coverage:"):
			if cov, err := extractCoverageFromString(row.Output); err == nil {
				v.Coverage = cov
				v.HasTestFiles = true
				v.HasTests = true
			}
		}

		results[row.Package] = v
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("got error from scanner: %w", err)
	}
	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("command did not finish successfully: %w", err)
	}
	return results, nil
}

type GoTestResultActionEnum string

const (
	GoTestResultActionPass   GoTestResultActionEnum = "pass"
	GoTestResultActionFail   GoTestResultActionEnum = "fail"
	GoTestResultActionSkip   GoTestResultActionEnum = "skip"
	GoTestResultActionOutput GoTestResultActionEnum = "output"
)

// GoPackageTestRunResultOutput is json formatted output result of a run go test
type GoPackageTestRunResultOutput struct {
	Action  GoTestResultActionEnum `json:"Action"`
	Package string                 `json:"Package"`
	Output  string                 `json:"Output"`
}

func extractCoverageFromString(input string) (float64, error) {
	for _, v := range strings.Split(input, " ") {
		if strings.HasSuffix(v, "%") {
			return strconv.ParseFloat(strings.TrimSuffix(v, "%"), 64)
		}
	}
	return 0, errors.New("not found")
}
