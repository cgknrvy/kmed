BIN_DIR := "bin"
APP_NAME := "kmed"
FRONTEND_DIR := "frontend"
DIST_DIR := "$(FRONTEND_DIR)/dist"

.PHONY: all build frontend linux windows run clean

all: build

frontend: 
	cd $(FRONTEND_DIR) && npm install && npm run build

# Build the frontend first so that it can be embeded in the go binary

# Builds for the current architecture.
build: frontend
	mkdir -p $(BIN_DIR)
	go build -tags="fts5" -0 $(BIN_DIR)/$(APP_NAME) .

# Builds for linux-am64
linux: frontend
	mkdir -p $(BIN_DIR)
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
							go build -tags="fts5" -o $(BIN_DIR)/$(APP_NAME)-linux-amd64 .

# Build for windows-amd64
windows: frontend
	mkdir -p $(BIN_DIR)
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 \
							CC=x86_64-w64-mingw32-gcc \
							go build -ldflags="-H=windowsgui" -tags="fts5" -o $(BIN_DIR)/$(APP_NAME)-windows-amd64.exe .

# Build the frontend before running so that the latest dist is embeded
run: frontend
	go run -tags="fts5" .


clean:
	rm -rf $(DIST_DIR)
	rm -f $(BIN_DIR)

