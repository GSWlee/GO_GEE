package gee

import (
	"log"
	"net/http"
)

type router struct {
	routers map[string]HandlerFunc
}

func NewRouter() *router {
	return &router{routers: map[string]HandlerFunc{}}
}

//添加路由器
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - $s",method,pattern)
	key := method + "-" + pattern
	r.routers[key] = handler
}

//执行路由
func (r *router) handle(c *Context)  {
	url:=c.Method+"-"+c.Path
	if handler,ok:=r.routers[url];ok{
		handler(c)
	}else {
		c.String(http.StatusNotFound,"404 NOT FOUND %s",c.Path)
	}
}