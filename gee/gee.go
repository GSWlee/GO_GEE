package gee

import (
	"net/http"
	"strings"
)

// HandlerFunc defines the request handler used by gee
type HandlerFunc func(c *Context)

//路由分组
type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc // support middleware
	parent      *RouterGroup  // support nesting
	engine      *Engine       // all groups share a Engine instance
}

// 路由表
type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup //存储所有的路由组
}

func (g *RouterGroup) Group(prefix string)  *RouterGroup{
	group:=&RouterGroup{
		prefix: prefix,
		engine: g.engine,
		parent: g,
	}
	g.engine.groups=append(g.engine.groups,group)
	return group
}

// New is the constructor of gee.Engine
func New() *Engine {
	engine:=&Engine{router: NewRouter()}
	engine.RouterGroup=&RouterGroup{engine: engine}
	engine.groups=[]*RouterGroup{engine.RouterGroup}
	return engine
}

// Run defines the method to start a http server
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

//添加路由
func (engine *Engine) addRoute(method string,pattern string,handle HandlerFunc)  {
	engine.router.addRoute(method,pattern,handle)
}

func (g *RouterGroup) addRoute(method string,comp string,handle HandlerFunc)  {
	pattern:=g.prefix+comp
	g.engine.router.addRoute(method,pattern,handle)
}

func (engine *Engine) GET(pattern string,handle HandlerFunc)  {
	engine.router.addRoute("GET",pattern,handle)
}

func (g *RouterGroup) GET(pattern string,handle HandlerFunc)  {
	g.addRoute("GET",pattern,handle)
}

func (engine *Engine) POST(pattern string,handle HandlerFunc)  {
	engine.router.addRoute("POST",pattern,handle)
}

func (g *RouterGroup) POST(pattern string,handle HandlerFunc)  {
	g.addRoute("POST",pattern,handle)
}

func (g *RouterGroup) Use(middlewares ...HandlerFunc)  {
	g.middlewares=append(g.middlewares,middlewares...)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middleware []HandlerFunc
	for _,group:=range engine.groups{
		if strings.HasPrefix(req.URL.Path,group.prefix){
			middleware=append(middleware,group.middlewares...)
		}
	}
	c:=NewContext(w,req)
	c.handlers=middleware
	engine.router.handle(c)
}