package web

import (
	"net/http"
	"strings"
)

// HandlerFunc defines the request handler used by ty
//type HandlerFunc func(http.ResponseWriter, *http.Request)
type HandlerFunc func(c *Context)


//defined route group
type RouterGroup struct {
	prefix string
	middlewares []HandlerFunc
	mind *Mind
}

//implement of ServerHttp
type Mind struct {
	router *router
	//group
	*RouterGroup //root RootGroup
	groups []*RouterGroup //all groups
}

//constructor of ty
func New() *Mind {
	mind := &Mind{router: newRouter()}      //new mind
	mind.RouterGroup = &RouterGroup{mind:mind} //new root routerGroup
	mind.groups = []*RouterGroup{mind.RouterGroup}

	return mind
	//return &Mind{router:newRouter()}
}

// Default use Recovery middlewares
func D() *Mind {
	m := New()
	m.Use(Recovery())
	return m
}

//add group
func (group *RouterGroup) Group(prefix string) * RouterGroup{
	mind := group.mind

	newGroup := &RouterGroup{
		prefix:      group.prefix + prefix,
		mind:        mind,
	}

	mind.groups = append(mind.groups, newGroup)
	return newGroup
}

//add middleware to routerGroup
func (group *RouterGroup) Use(handlerFuncs ...HandlerFunc)  {
	group.middlewares = append(group.middlewares, handlerFuncs...)
}

//add route to routerGroup
func (group *RouterGroup) addRoute(method string, pattern string, handler HandlerFunc) {
	pattern = group.prefix + pattern
	group.mind.router.addRoute(method, pattern, handler)
}

//add GET request route to routerGroup
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

//add POST request route to routerGroup
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

//run a HTTP server
func (m *Mind) Run(addr string) error {
	return http.ListenAndServe(addr, m)
}

//handle request
func (m *Mind) ServeHTTP (w http.ResponseWriter, r *http.Request) {
	middlewares := make([]HandlerFunc,0)
	for _, group := range m.groups {
		//find middlewares by group
		if strings.HasPrefix(r.URL.Path, group.prefix){
			middlewares = append(middlewares,group.middlewares...)
		}
	}
	//init context
	c := newContext(w,r)
	//add middleware to context
	c.handlers = middlewares

	m.router.Handle(c)
}

