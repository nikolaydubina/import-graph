clean:
	-rm -rf bin

build: clean
	GO111MODULE=on go build -o bin/import-graph app/import-graph/main.go
	GO111MODULE=on go build -o bin/jsonl-to-dot app/jsonl-to-dot/main.go

docs: build
	cd third-party; git submodule update --init --recursive

	cd third-party/go-featureprocessing; go mod graph > gomodgraph
	#cat third-party/go-featureprocessing/gomodgraph | ./bin/import-graph -output=dot | dot -Tsvg > docs/go-featureprocessing/output.dot.svg

	cd third-party/calendarheatmap; go mod graph > gomodgraph
	cat third-party/calendarheatmap/gomodgraph | ./bin/import-graph > docs/calendarheatmap/output.jsonl
	cat third-party/calendarheatmap/gomodgraph | ./bin/import-graph -output=dot > docs/calendarheatmap/output.dot
	cd docs/calendarheatmap; cat output.dot | dot -Tsvg > output.dot.svg

.PHONY: clean build docs
