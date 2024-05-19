package main

import (
	"io"
	"os"
	"social/project/api/router"
	"social/project/database"
	"social/project/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetUpLogOutput() {
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
}
func main() {
	SetUpLogOutput()
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(gin.Recovery(), middleware.Logger())
	r.Use(cors.Default())
	// log.Printf(os.Getenv("DATABASE_URL"))
	database.Connect()

	router.InitializeUserRouter(r)

	r.GET("/", middleware.AuthMiddleware(), func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Homepage",
		})
	})

	r.Run(":5000")
}
