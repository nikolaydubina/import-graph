> When was the last time my dependencies were updated?  
> Are they well maintained?  
> Who maintains them?  
> Who else uses them?  
> Do they have good test coverage?  

Tools here allow you to _collect_ and _visualize_ data about your dependencies.
Every stage is a UNIX filter which allows great integrations.
Front ends are in [dot](https://graphviz.org) and web.
Currently works with `Go` and `GitHub`.

## Related Projects

- `Graphviz` https://graphviz.org/ is a very popular tool for visualizing graph data, most of tools bellow use dot from it
- `Graphviz` https://graphviz.org/Gallery/directed/neural-network.html is nice example of dot format
- `Docs` https://awesomeopensource.com/projects/dependency-graph is a list of dependency visualization projects  
- `Go` https://github.com/loov/goda written in Go; analyses imports on its own; does not collect dta; CLI; dot  
- `Go` https://github.com/adonovan/spaghetti wirtten in Go; search and read details about selected package; web; not graphic
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
