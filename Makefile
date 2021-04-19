clean:
	-rm -rf bin

build: clean
	GO111MODULE=on go build -o bin/import-graph main.go

.PHONY: clean build
