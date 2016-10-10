GO=go
BIN=./ahlt
SRC=src/ahlt/main.go
GOPATH=$(shell pwd)

all: clean golib build install

clean: 
	$(RM) $(BIN)
	$(RM) -rf ./src/github.com
	$(RM) /bin/wrk
	$(RM) -rf ./wrk	
wrk:
	@$(SHELL) -c "git clone https://github.com/wg/wrk.git"
	@$(SHELL) -c "cd wrk && sudo make && sudo install -m 0755 wrk /bin"

golib:
	@echo "Download Golang Library"
	@GOPATH=$(GOPATH) $(GO) get -u github.com/kataras/iris/iris
	@GOPATH=$(GOPATH) $(GO) get github.com/googollee/go-socket.io
	@GOPATH=$(GOPATH) $(GO) get github.com/PuerkitoBio/goquery
	@GOPATH=$(GOPATH) $(GO) get github.com/mattn/go-sqlite3
	@GOPATH=$(GOPATH) $(GO) get -u github.com/jinzhu/gorm
	@GOPATH=$(GOPATH) $(GO) get -u github.com/kataras/go-template
	@GOPATH=$(GOPATH) $(GO) get -u github.com/flosch/pongo2

build: $(SRC)
	@echo "Build"
	@GOPATH=$(GOPATH) $(GO) build -o $(BIN) $+

install: 
	sudo install -m 0755 ahlt /bin

