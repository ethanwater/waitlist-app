package main

import (
	"net/mail"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

type App struct {
	Address string
	Engine  *gin.Engine
	API     *gin.RouterGroup
}

func main() {
	router := gin.Default()
	app := App{
		":8080",
		router,
		router.Group("/api"),
	}

	app.API.GET("/enrollemail", enroll())

	app.Engine.Use(static.Serve("/", static.LocalFile("./views", true)))
	app.Engine.Run(app.Address)
}

func enroll() func(*gin.Context) {
	return func(ctx *gin.Context) {
		email := ctx.Query("email")

		//validate email
		_, err := mail.ParseAddress(email)
		if err != nil {
			ctx.JSON(200, gin.H{
				"error": "invalid email",
			})
			return
		}

		//store email in waitlist

	}
}
