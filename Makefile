GO=go
BIN=./ahlt
SRC=src/ahlt/main.go
GOPATH=$(shell pwd)

all: clean depen install

clean: 
	$(RM) $(BIN)
	$(RM) -rf ./src/github.com
	$(RM) /bin/wrk
	$(RM) -rf ./wrk	
depen:
	@$(SHELL) -c "git clone https://github.com/wg/wrk.git"
	@$(SHELL) -c "cd wrk && sudo make"
	@$(SHELL) -c "cd wrk && sudo cp wrk /bin/"
	@echo "Download Golang Library"
	@GOPATH=$(GOPATH) $(GO) get -u github.com/kataras/iris/iris
	@GOPATH=$(GOPATH) $(GO) get github.com/googollee/go-socket.io
	@GOPATH=$(GOPATH) $(GO) get github.com/PuerkitoBio/goquery
	@GOPATH=$(GOPATH) $(GO) get github.com/mattn/go-sqlite3
	@GOPATH=$(GOPATH) $(GO) get -u github.com/jinzhu/gorm
	@GOPATH=$(GOPATH) $(GO) get -u github.com/kataras/go-template
	@GOPATH=$(GOPATH) $(GO) get -u github.com/flosch/pongo2

install: $(SRC)
	@echo "Build"
	@GOPATH=$(GOPATH) $(GO) build -o $(BIN) $+

