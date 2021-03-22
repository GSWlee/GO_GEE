package gee

import (
	"html/template"
	"net/http"
	"path"
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
	htmlTemplate *template.Template
	funcmap template.FuncMap
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

// create static handler
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		// Check if file exists and/or if we have permission to access it
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		fileServer.ServeHTTP(c.Writer, c.Request)
	}
}

// serve static files
func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	// Register GET handlers
	//在注册路由时，完成url路径对于服务器上实际文件系统的映射
	//返回的handle中的fileServer对应的时给定的映射
	group.GET(urlPattern, handler)
}

func (e *Engine)SetFuncmap(funcmap template.FuncMap)  {
	e.funcmap=funcmap
}

func (e *Engine)LoadHTMLGlob(pattern string)  {
	e.htmlTemplate=template.Must(template.New("").Funcs(e.funcmap).Parse(pattern))
}
