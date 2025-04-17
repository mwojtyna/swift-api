CSV=swift-codes.csv
BIN_DIR=bin
SERVER_BIN=$(BIN_DIR)/server
PARSER_BIN=$(BIN_DIR)/parse-swift

build-server: 
	@go build -o $(SERVER_BIN) cmd/server/main.go

build-parser: 
	@go build -o $(PARSER_BIN) cmd/parser/main.go

serve: build-server
	@./$(SERVER_BIN)

parse: build-parser
	@./$(PARSER_BIN) $(CSV)

test:
	@go test ./...

clean:
	@rm -rf bin

.PHONY: build-server build-parser serve parse test clean
