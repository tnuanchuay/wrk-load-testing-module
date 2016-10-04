GO = go
BIN = ./ahlt
SRC = src/ahlt/main.go
OLD = ""
all: clean install

clean: 
	$(RM) $(BIN)

install: $(SRC)
	$(GO) build -o $(BIN) $+

