package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	//基本
	Writer http.ResponseWriter
	Request *http.Request

	Method string
	Path string
	Params map[string]string
	StatusCode int

	//middleware
	handlers []HandlerFunc
	index int

	engine *Engine
}

func NewContext(w http.ResponseWriter,req * http.Request) *Context {
	return &Context{
		Writer: w,
		Request: req,
		Method: req.Method,
		Path: req.URL.Path,
		index: -1,
	}
}

func (c *Context) PostForm(key string) string {
	return c.Request.FormValue(key)
}

func (c *Context) Query(key string) string{
	return c.Request.URL.Query().Get(key)
}

func (c *Context) Status(code int)  {
	c.StatusCode=code
	c.Writer.WriteHeader(code)
}

func (c *Context) SetStatus(key string,value string)  {
	c.Writer.Header().Set(key,value)
}

func (c *Context) String(code int,format string,values ...interface{})  {
	c.SetStatus("Content-Type","text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format,values...)))
}

func (c *Context) JSON(code int,obj interface{})  {
	c.SetStatus("Content-Type","application/json")
	c.Status(code)
	encode:=json.NewEncoder(c.Writer)
	if err:=encode.Encode(obj);err!=nil{
		http.Error(c.Writer,err.Error(),code)
	}
}

func (c *Context) Fail(code int,err string)  {
	http.Error(c.Writer,err,code)
}

func (c *Context) Data(code int ,obj []byte)  {
	c.Status(code)
	c.Writer.Write(obj)
}

func (c *Context) HTML(code int,name string,data interface{})  {
	c.Status(code)
	c.SetStatus("Content-Type","text/html")
	if err:=c.engine.htmlTemplate.ExecuteTemplate(c.Writer,name,data);err!=nil{
		c.Fail(http.StatusInternalServerError,err.Error())
	}
}

func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

func (c *Context) Next()  {
	c.index++
	s :=len(c.handlers)
	for ;c.index<s;c.index++{
		c.handlers[c.index](c)
	}
}