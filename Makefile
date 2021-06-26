test:
	go generate ./...
	go test -covermode=atomic ./...

build: test
	go build

docs-run: 
	cat docs/gin.gomodgraph | ./import-graph -i=gomod > docs/gin.jsonl

docs-render:
	cat docs/gin.jsonl | jsonl-graph -color-scheme=file://$$PWD/basic.json | dot -Tsvg > docs/gin.svg

docs-gin-gomodgraph:
	git clone https://github.com/gin-gonic/gin
	cd gin; go mod graph > ../docs/gin.gomodgraph
	rm -rf gin

.PHONY: build docs-run docs-render docs-gin-gomodgraph
