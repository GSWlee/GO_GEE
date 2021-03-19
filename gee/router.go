package gee

import (
	"log"
	"net/http"
	"strings"
)

type node struct {
	pattern string     //待匹配路由
	part string        //节点值
	children []*node   //子节点列表
	isWild bool        //是否为模糊查找，带有*，：等
}

type router struct {
	roots map[string]*node
	routers map[string]HandlerFunc
}

//查找第一个匹配的节点
func (n *node) mathChild(part string) *node {
	for _,item:=range n.children{
		if item.part==part||item.isWild{
			return item
		}
	}
	return nil
}

//查找所有匹配节点
func (n *node) mathChildren(part string)  []*node {
	children:=make([]*node,0)
	for _,item:=range n.children{
		if item.part==part||item.isWild{
			children=append(children,item)
		}
	}
	return children
}

//插入路由
func (n *node) insert(pattern string,parts []string,height int)  {
	if len(parts)==height{
		n.pattern=pattern
		return
	}

	part:=parts[height]
	child:=n.mathChild(part)
	if child==nil{
		child=&node{part: part,isWild: part[0]=='*'||part[0]==':'}
		n.children=append(n.children,child)
	}
	child.insert(pattern,parts,height+1)
}

//查找路由
func (n *node) search(parts []string,height int) *node{
	if len(parts)==height||strings.HasPrefix(n.part,"*"){
		if n.pattern==""{
			return nil
		}
		return n
	}

	part:=parts[height]
	children:=n.mathChildren(part)
	for _,child:=range children{
		result:=child.search(parts,height+1)
		if result!=nil{
			return result
		}
	}

	return nil
}

func NewRouter() *router {
	return &router{
		routers: map[string]HandlerFunc{},
		roots: map[string]*node{},
	}
}

func splitPattern(pattern string) []string {
	vs:=strings.Split(pattern,"/")
	parts:=[]string{}
	for _,part:=range vs{
		if part!=""{
			parts=append(parts,part)
			if part[0]=='*'{
				break
			}
		}
	}
	return parts
}

//添加路由器
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - " +"%s",method,pattern)
	parts:=splitPattern(pattern)
	key := method + "-" + pattern
	if _,ok:= r.roots[method];!ok{
		r.roots[method]=&node{}
	}
	r.roots[method].insert(pattern,parts,0)
	r.routers[key] = handler
}

//选择路由go
func (r *router) getRoute(method string,path string) (*node, map[string]string) {
	params:=map[string]string{}
	searchparts:=splitPattern(path)
	root,ok:=r.roots[method]
	if !ok{
		return nil,nil
	}
	n:=root.search(searchparts,0)
	if n!=nil{
		parts:=splitPattern(n.pattern)
		for index,value:=range parts{
			if value[0]==':'{
				params[value[1:]]=searchparts[index]
			}
			if value[0]=='*' {
				params[value[1:]] = strings.Join(searchparts[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

//执行路由
func (r* router) handle(c *Context)  {
	n,parms:=r.getRoute(c.Method,c.Path)
	if n!=nil{
		c.Params=parms
		log.Println(c)
		key:=c.Method+"-"+n.pattern
		c.handlers=append(c.handlers,r.routers[key])
	}else{
		c.String(http.StatusNotFound,"404 NOT FOUND: %s\n",c.Path)
	}
	c.Next()
}