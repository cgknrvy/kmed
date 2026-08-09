BIN_DIR := "bin"
APP_NAME := "kmed"
APP_OUTPUT := "$(BIN_DIR)/$(APP_NAME)"
FRONTEND_DIR := "frontend"
DIST_DIR := "$(FRONTEND_DIR)/dist"

.PHONY: all build frontend backend run clean

all: build

# Build the frontend first so that it can be embeded in the go binary
build: frontend backend

frontend: 
	cd $(FRONTEND_DIR) && npm install && npm run build

backend:
	go build -tags="fts5" -o $(APP_OUTPUT) .

run:
	go run -tags="fts5" .


clean:
	rm -rf $(DIST_DIR)
	rm -f $(APP_OUTPUT)

