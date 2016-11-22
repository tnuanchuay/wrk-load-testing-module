GO=go
BIN=./ahlt
SRC=src/ahlt/main.go
GOPATH=$(shell pwd)

all: golib build

clean: 
	$(RM) $(BIN)

golib:
	@echo "Download Golang Library"
	@GOPATH=$(GOPATH) $(GO) get -u github.com/kataras/iris/iris
	@GOPATH=$(GOPATH) $(GO) get github.com/googollee/go-socket.io
	@GOPATH=$(GOPATH) $(GO) get github.com/PuerkitoBio/goquery
    	@GOPATH=$(GOPATH) $(GO) get github.com/mattn/go-sqlite3
	@GOPATH=$(GOPATH) $(GO) get -u github.com/jinzhu/gorm
	@GOPATH=$(GOPATH) $(GO) get -u github.com/flosch/pongo2

build: $(SRC)
	@echo "Build"
	@GOPATH=$(GOPATH) $(GO) build -o $(BIN) $+

