.PHONY: build build-server build-cli test fmt clean

build: build-server build-cli

build-server:
	go build -o azlake ./cmd/azlake

build-cli:
	go build -o azlakectl ./cmd/azlakectl

test:
	go test -count=1 -race ./...

fmt:
	gofmt -w .

clean:
	rm -f azlake azlakectl
