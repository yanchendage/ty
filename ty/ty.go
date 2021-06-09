package ty

import (
	"net/http"
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
	mind := &Mind{router:newRouter()}//new mind
	mind.RouterGroup = &RouterGroup{mind:mind}//new root routerGroup
	mind.groups = []*RouterGroup{mind.RouterGroup}

	return mind
	//return &Mind{router:newRouter()}
}

//add group
func (group *RouterGroup) Group(prefix string) * RouterGroup{
	mind := group.mind

	newGroup := &RouterGroup{
		prefix:      group.prefix + prefix,
		parent:      group,
		mind:        mind,
	}

	mind.groups = append(mind.groups, newGroup)
	return newGroup
}

//router
func (group *RouterGroup) addRoute(method string, pattern string, handler HandlerFunc) {
	pattern = group.prefix + pattern
	group.mind.router.addRoute(method, pattern, handler)
}

//add GET request
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

//add POST request
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

//run a HTTP server
func (m *Mind) Run(addr string) error {
	return http.ListenAndServe(addr, m)
}

//handle request
func (m *Mind) ServeHTTP (w http.ResponseWriter, r *http.Request) {
	//new Context
	c := newContext(w,r)
	m.router.Handle(c)
}