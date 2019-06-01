package main

import (
	"github.com/asphaltbot/file-storage/routes"
	"github.com/asphaltbot/file-storage/util"
	"github.com/getsentry/raven-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/sentry"
	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/acme/autocert"
	"log"
	"time"
)

func main() {
	r := gin.Default()

	r.Use(sentry.Recovery(raven.DefaultClient, false))
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods: []string{"GET", "POST", "DELETE"},
		AllowHeaders: []string{"Origin"},
		ExposeHeaders: []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge: 24 * time.Hour,
	}))

	registerRoutes(r)

	r.GET("/", Index)

	if util.IsRunningInProd() {
		gin.SetMode(gin.ReleaseMode)

		m := autocert.Manager{
			Prompt: autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist("storage.asphaltbot.com"),
			Cache: autocert.DirCache("/home/storage/cache"),
		}

		log.Fatal(autotls.RunWithManager(r, &m))

	} else {
		gin.SetMode(gin.DebugMode)
		_ = r.Run(":80")
	}

}

func Index(c *gin.Context) {
	c.Data(200, "text/html", []byte("<!DOCTYPE html><html lang=\"en\" dir=\"ltr\"> <head> <meta charset=\"utf-8\"> <title>Asphalt Storage Server</title> <link href=\"https://fonts.googleapis.com/css?family=Roboto&display=swap\" rel=\"stylesheet\"> <style>h1, p, small{font-family: 'Roboto', sans-serif;}</style> </head> <body> <div align=\"center\"> <h1>Asphalt Storage Server</h1> <p>This server is running custom open source backup software, click <a href=\"https://github.com/asphaltbot/file-storage\">here</a> for more information.</p><br/> <small>Created by Connor Wright</small> </div></body></html>"))
}


func registerRoutes(e *gin.Engine) {
	routes.RegisterUploadRoutes(e)
}