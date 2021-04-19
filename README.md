> Are my dependencies any good?

Collect data about your dependencies, visualize with [dot](https://graphviz.org) or with `jsonl-graph-viz` (TODO: git, website).
This is powering (TODO: website all together).

## Examples

```
$ go mod graph | ./bin/import-graph -output=dot | dot -Tsvg > output.svg
```
TODO: image

```
$ go mod graph | ./bin/import-graph | jq
...
{
    "id": "github.com/google/gofuzz",
    "module_name": "github.com/google/gofuzz",
    "last_commit": "2021-01-06T07:40:07+08:00",
    "days_since_last_commit": 102.72333373628472,
    "years_since_last_commit": 0.3057242075485009,
    "months_since_last_commit": 3.668690490582011,
    "num_contributors": 19,
    "has_tests": true,
    "has_test_files": true,
    "num_packages": 2,
    "num_packages_with_tests": 2,
    "num_packages_with_tests_files": 2,
    "num_packages_tests_passed": 2,
    "min_package_coverage": 85.7,
    "avg_package_coverage": 86.95
}
...
{
  "from": "github.com/nikolaydubina/go-featureprocessing",
  "to": "github.com/google/gofuzz"
}
```

## Data Sources

- [x] Runs tests and code coverage
- [x] Analyzes git log
- [ ] Runs linters
- [ ] Detects benchmarks
- [ ] Checks for Codecov Codacy
- [ ] Checks GoReportCard
- [ ] Checks Readme for words and badges
- [ ] Version is stable (same as godoc)
- [ ] Mentioning in Awesome-lists
- [ ] Analyzes GitHub Repo, Organization, Contributors, Stars, Contributor, Contributors organizations, Is contributor active?

## Related Projects

- `Graphviz` https://graphviz.org/ is a very popular tool for visualizing graph data, most of tools bellow use dot from it
- `Graphviz` https://graphviz.org/Gallery/directed/neural-network.html is nice example of dot format
- `Graphviz` http://magjac.com/graphviz-visual-editor online editor
- `Docs` https://awesomeopensource.com/projects/dependency-graph is a list of dependency visualization projects  
- `Go` https://github.com/loov/goda written in Go; analyses imports on its own; does not collect dta; CLI; dot  
- `Go` https://github.com/adonovan/spaghetti wirtten in Go; search and read details about selected package; web; not graphic
- `Go` https://github.com/psampaz/go-mod-outdated written in Go; pretty prints built-in Go tooling; not collecting; not dot; CLI 
- `Go` https://github.com/firstrow/go-outdated written in Go; deprecated; collects data from GitHub; CLI
- `GitHub Tasks` https://github.com/moul/depviz written in Go and JavaScript; collects data from GitHub Tasks; plots nice web UI
- `C++` https://github.com/jmarkowski/codeviz written in Python; C++ headers analysis; does not collect data; CLI; dot  
- `Python` https://github.com/thebjorn/pydeps written in Python; looks for Python bytecode imports; clustering; does not collect data; CLI; dot  
- `Python` https://github.com/naiquevin/pipdeptree written in Python; looks for python modules locally; does not collect data; CLI; JSON and dot, Deprecated  
- `JavaScript` https://github.com/auchenberg/dependo written in JavaScript; does not fetch data; D3.js; CLI; HTML   
- `JavaScript` https://github.com/pahen/madge written in JavaScript; does not collect data; CLI; dot  
- `JavaScript` https://github.com/sverweij/dependency-cruiser written in JavaScript; rules; does not collect data; CLI; dot  
- `JavaScript` https://github.com/anvaka/npmgraph.an written in JavaScript; collects data; HTML; hosted in GitHub Pages  
- `JavaScript` https://github.com/anvaka/npmgraphbuilder written in JavaScript; collects data; module  
- `JavaScript` https://github.com/dyatko/arkit written in JavaScript; modules and dependencies; CLI; svg, puml  
- `JavaScript` https://github.com/hughsk/colony written in JavaScript; does not collect data; HTML; JSON  
- `JavaScript` https://www.npmjs.com/package/node-dependency-visualizer written in JavaScript; does not collectdata; CLI; dot  
- `Objective-C` `Swift` https://github.com/PaulTaykalo/objc-dependency-visualizer written in JavaScript and Ruby; does not collect data; CLI; dot; HTML; D3.js   
- `Java` https://github.com/arunkumar9t2/scabbard written in Kotlin; CLI; dot  
- `PHP` https://github.com/mamuz/PhpDependencyAnalysis written in PHP; does not collect data; code analysis; CLI; dot  
- `Go` `Python` `Java` `JavaScript` `C++` https://github.com/oss-review-toolkit/ort written in Kotlin JavaSCript Python; collects data; analyses; analysis, downloading, reporting; used for licence scanning in open source; good architecture; a bit lacking support for Go; components may not be used separately  
- `Code` https://github.com/aspiers/git-deps written in Python; analyses dependencies of commits in Git repository  
