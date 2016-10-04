GO=go
BIN=./ahlt
SRC=src/ahlt/main.go
GOPATH=$(shell pwd)

all: clean depen install

clean: 
	$(RM) $(BIN)
	$(shell git reset --hard)
depen:
	#git  clone https://github.com/wg/wrk.git
	#$(shell cd wrk; make;)
	#$(shell ln -s wrk /bin/)
	@GOPATH=$(GOPATH) $(GO) get -u github.com/kataras/iris/iris
	@GOPATH=$(GOPATH) $(GO) get github.com/googollee/go-socket.io
	@GOPATH=$(GOPATH) $(GO) get github.com/PuerkitoBio/goquery
	@GOPATH=$(GOPATH) $(GO) get github.com/mattn/go-sqlite3
	@GOPATH=$(GOPATH) $(GO) get -u github.com/jinzhu/gorm
	@GOPATH=$(GOPATH) $(GO) get -u github.com/kataras/go-template
	@GOPATH=$(GOPATH) $(GO) get -u github.com/flosch/pongo2


install: $(SRC)
	@GOPATH=$(GOPATH) $(GO) build -o $(BIN) $+

