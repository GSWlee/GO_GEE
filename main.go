package main

import (
	"./gee"
	"log"
	"net/http"
	"time"
)

func onlyForV2() gee.HandlerFunc {
	return func(c *gee.Context) {
		// Start timer
		t := time.Now()
		c.Next()
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Request.RequestURI, time.Since(t))
	}
}

func main()  {


	//新建web框架
	app:=gee.New()

	//配置路由
	app.Use(gee.Logger())
	app.GET("/",indexHandler)
	app.GET("/hello",helloHandler)
	app.POST("/login",loginhandler)
	api:=app.Group("/api")
	api.GET("/name", func(c *gee.Context) {
		c.String(http.StatusOK,"name is qqq\n")
	})
	api.Use(onlyForV2())
	app.GET("/hello/:name", func(c *gee.Context) {
		c.String(http.StatusOK,"hello %s you're at %s\n",c.Param("name"),c.Path)
	})
	app.GET("/assets/*filepath", func(c *gee.Context) {
		c.JSON(http.StatusOK,gee.H{"filepath":c.Param("filepath")})
	})

	//执行框架
	log.Fatal(app.Run(":8080"))
}

// show r.url
func indexHandler(c *gee.Context)  {
	c.HTML(http.StatusOK,"<h1>hello world</h1>")
}

//show r.header
func helloHandler(c *gee.Context)  {
	c.String(http.StatusOK,"hello %s you are at %s\n",c.Query("name"),c.Path)
}

//login
func loginhandler(c *gee.Context)  {
	c.JSON(http.StatusOK,gee.H{
		"username":c.PostForm("username"),
		"password":c.PostForm("password"),
	})
}

