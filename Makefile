build-linux-amd64:
	GOOS=linux GOARCH=amd64 go build -o ./bin/linux-amd64/server ./cmd/server
	GOOS=linux GOARCH=amd64 go build -o ./bin/linux-amd64/epoxy ./cmd/epoxy

build-darwin-arm64:
	GOOS=darwin GOARCH=arm64 go build -o ./bin/darwin-arm64/server ./cmd/server
	GOOS=darwin GOARCH=arm64 go build -o ./bin/darwin-arm64/epoxy ./cmd/epoxy

build-windows-amd64:
	GOOS=windows GOARCH=amd64 go build -o ./bin/windows-amd64/server.exe ./cmd/server
	GOOS=windows GOARCH=amd64 go build -o ./bin/windows-amd64/epoxy.exe ./cmd/epoxy

clean:
	rm -rf ./bin/*

all: clean build-linux-amd64 build-darwin-arm64 build-windows-amd64
