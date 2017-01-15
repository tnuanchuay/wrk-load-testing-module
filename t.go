package main

import (
	"github.com/kataras/iris"
	"time"
)

func main(){
	iris.Get("/", func(ctx *iris.Context){
		time.Sleep(1 * time.Millisecond)
		ctx.JSON(200, "ok")
	})
	iris.Listen(":8080")
}
