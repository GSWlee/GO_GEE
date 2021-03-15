package gee

import (
	"net/http"
)

// HandlerFunc defines the request handler used by gee
type HandlerFunc func(c *Context)

// 路由表
type Engine struct {
	router *router

}

// New is the constructor of gee.Engine
func New() *Engine {
	return &Engine{router: NewRouter()}
}

// Run defines the method to start a http server
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

//添加路由
func (engine *Engine) addRoute(method string,pattern string,handle HandlerFunc)  {
	engine.router.addRoute(method,pattern,handle)
}

func (engine *Engine) GET(pattern string,handle HandlerFunc)  {
	engine.router.addRoute("GET",pattern,handle)
}

func (engine *Engine) POST(pattern string,handle HandlerFunc)  {
	engine.router.addRoute("POST",pattern,handle)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c:=NewContext(w,req)
	engine.router.handle(c)
}