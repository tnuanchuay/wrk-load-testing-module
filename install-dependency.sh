#!/usr/bin/env bash
export GOPATH=$(pwd)
echo 'Getting iris...'
go get -u github.com/kataras/iris/iris
echo 'Getting go-socket.io'
go get github.com/googollee/go-socket.io
echo 'Getting goquery'
go get github.com/PuerkitoBio/goquery
echo 'Getting sqlite'
go get github.com/mattn/go-sqlite3
echo 'Getting gorm'
go get -u github.com/jinzhu/gorm
echo 'Getting iris/Go-Template'
go get -u github.com/kataras/go-template
