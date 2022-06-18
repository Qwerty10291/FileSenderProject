package main

import "github.com/gin-gonic/gin"

func main(){
	app := gin.Default()
	app.GET("/", func(ctx *gin.Context) {
		ctx.String(200, "Hello, file sender!")
	})
}