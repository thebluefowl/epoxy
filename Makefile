# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get
BIN_DIR = /bin

# List of binaries to build
BINARY_NAMES = server epoxyd

all: test build
.PHONY: all

build: $(BINARY_NAMES)
.PHONY: build

$(BINARY_NAMES):
	$(GOBUILD) -o $(BIN_DIR)/$@ ./cmd/$@

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	for binary in $(BINARY_NAMES); do \
		rm -f $(BIN_DIR)/$$binary; \
	done

run: build
	@echo "Specify the binary name to run. Example: make run BINARY_NAME=server"

deps:
	$(GOGET) -v ./...

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v
