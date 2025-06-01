package main

import "github.com/gin-gonic/gin"

func run() {
	r := gin.Default()

	r.StaticFile("/favicon.ico", "./static/favicon.ico")
	r.Static("/static", "./static")
	r.StaticFile("/", "./html/index.html")
	r.StaticFile("/invitation", "./html/invitation.html")

	apiGroup := r.Group("/api")
	{
		apiGroup.POST("genSTSCreds", genSTSCredsHandler)
		apiGroup.GET("getObject", getObjectHandler)
		apiGroup.POST("putObject", putObjectHandler)
	}

	r.Run(":8080")
}
