clean:
	-rm -rf bin

build: clean
	GO111MODULE=on go build -o bin/import-graph app/import-graph/main.go
	GO111MODULE=on go build -o bin/jsonl-to-dot app/jsonl-to-dot/main.go

docs: build
	cd third-party; git submodule update --init --recursive
	-rm -rf docs/gin_*
	cd third-party/gin; go mod graph > ../../docs/gin_graph
	cat docs/gin_graph | ./bin/import-graph -test > docs/gin.jsonl
	cat docs/gin.jsonl | ./bin/jsonl-to-dot > docs/gin.dot
	cat docs/gin.dot | dot -Tsvg > docs/gin.svg

.PHONY: clean build docs
