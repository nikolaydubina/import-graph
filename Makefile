clean:
	-rm -rf bin

install-tools:
	go get -u honnef.co/go/tools/cmd/staticcheck

build: clean
	GO111MODULE=on go build -o bin/import-graph app/import-graph/main.go
	GO111MODULE=on go build -o bin/jsonl-to-dot app/jsonl-to-dot/main.go

lint:
	go vet ./...
	go fmt ./...
	staticcheck ./...

test:
	go generate ./...
	go test -covermode=atomic ./...

docs: build
	-rm docs/gin_*
	-rm -rf docs/gin
	cd docs; git clone https://github.com/gin-gonic/gin
	cd docs/gin; go mod graph > ../../docs/gin_graph
	cat docs/gin_graph | ./bin/import-graph -test > docs/gin.jsonl
	cat docs/gin.jsonl | ./bin/jsonl-to-dot > docs/gin.dot
	cat docs/gin.dot | dot -Tsvg > docs/gin.svg
	cat docs/gin.jsonl | ./bin/jsonl-to-dot -color-scheme=file://$$PWD/app/import-graph/basic-colors.json > docs/gin_color.dot
	cat docs/gin_color.dot | dot -Tsvg > docs/gin_color.svg
	-rm -rf docs/gin

docs-render:
	cat docs/gin.jsonl | ./bin/jsonl-to-dot > docs/gin.dot
	cat docs/gin.dot | dot -Tsvg > docs/gin.svg
	cat docs/gin.jsonl | ./bin/jsonl-to-dot -color-scheme=file://$$PWD/app/import-graph/basic-colors.json > docs/gin_color.dot
	cat docs/gin_color.dot | dot -Tsvg > docs/gin_color.svg

.PHONY: clean build docs
