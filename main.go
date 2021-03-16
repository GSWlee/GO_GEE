package main

import (
	"./gee"
	"log"
	"net/http"
)

func main()  {


	//新建web框架
	app:=gee.New()

	//配置路由
	app.GET("/",indexHandler)
	app.GET("/hello",helloHandler)
	app.POST("/login",loginhandler)
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

