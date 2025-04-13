CSV=swift-codes.csv
BIN_DIR=bin
SERVER_BIN=$(BIN_DIR)/server
PARSER_BIN=$(BIN_DIR)/parse-swift

build: 
	@go build -o $(SERVER_BIN) cmd/server/main.go

build-parser: 
	@go build -o $(PARSER_BIN) cmd/parser/main.go

serve: build
	@./$(SERVER_BIN)

parse: build-parser
	@./$(PARSER_BIN) $(CSV)

clean:
	@rm -rf bin

.PHONY: build build-parser serve parse clean
