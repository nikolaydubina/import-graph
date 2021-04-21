clean:
	-rm -rf bin

build: clean
	GO111MODULE=on go build -o bin/import-graph app/import-graph/main.go
	GO111MODULE=on go build -o bin/jsonl-to-dot app/jsonl-to-dot/main.go

docs-init:
	cd third-party; git submodule update --init --recursive

docs-calendarheatmap:
	cd third-party/calendarheatmap; go mod graph > gomodgraph
	cat third-party/calendarheatmap/gomodgraph | ./bin/import-graph -test > docs/calendarheatmap/output.jsonl
	cat third-party/calendarheatmap/gomodgraph | ./bin/import-graph -test -output=dot > docs/calendarheatmap/output.dot
	cd docs/calendarheatmap; cat output.dot | dot -Tsvg > output.dot.svg
	-rm third-party/calendarheatmap/gomodgraph

docs-go-featureprocessing:
	cd third-party/go-featureprocessing; go mod graph > gomodgraph
	cat third-party/go-featureprocessing/gomodgraph | ./bin/import-graph -output=dot | dot -Tsvg > docs/go-featureprocessing/output.dot.svg
	-rm third-party/go-featureprocessing/gomodgraph

docs: build docs-init docs-calendarheatmap

.PHONY: clean build docs
