SHELL=/bin/sh
GOCMD=go
GOBUILD=$(GOCMD) build
GOBIN=$(GOPATH)/bin

SERVER_DIR=server
CLIENT_DIR=client
SERVER_BINARY=messagingserver
CLIENT_BINARY=messagingclient

build: | proto
	$(GOBUILD) -o $(GOBIN)/$(SERVER_BINARY) $(SERVER_DIR)/cmd/main.go
	$(GOBUILD) -o $(GOBIN)/$(CLIENT_BINARY) $(CLIENT_DIR)/cmd/main.go
proto:
	pushd $(SERVER_DIR) && protoc --go_out=plugins=grpc:pkg api/*.proto
clean:
	go clean
	rm -f $(GOBIN)/$(SERVER_BINARY)
	rm -f $(GOBIN)/$(CLIENT_BINARY)
	rm -rf $(SERVER_DIR)/pkg/apis