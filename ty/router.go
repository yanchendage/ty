package ty

import (
	"net/http"
	"strings"
)

type router struct {
	roots map[string]*trie
	handlers map[string]HandlerFunc
}

func newRouter() *router{

	return &router{
		roots: make(map[string]*trie),
		handlers:make(map[string]HandlerFunc),
	}
}

func parsePattern(pattern string) []string{
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			//append only first *
			if item[0] == '*' {
				break
			}
		}
	}

	return parts
}

func (r *router) addRoute (method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)

	key := method + "-" + pattern

	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &trie{}
	}

	r.roots[method].add(pattern, parts, 0)
	r.handlers[key] = handler
}

func (r *router) findRoute(method string, path string)(*trie,map[string]string) {

	pathParts := parsePattern(path)
	params := make(map[string]string)

	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}

	n := root.find(pathParts, 0)

	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = pathParts[index]
			}

			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(pathParts[index:], "/")
				break
			}
		}
		return n, params
	}

	return  nil, nil
}

func (r *router) GET(pattern string, handler HandlerFunc) {
	r.addRoute("GET", pattern, handler)
}

func (r *router) POST(pattern string, handler HandlerFunc) {
	r.addRoute("POST", pattern, handler)
}

func (r *router) Handle(c *Context)  {

	n, param := r.findRoute(c.Method, c.Path)

	if n != nil {
		c.Params = param
		key := c.Method + "-" + n.pattern
		r.handlers[key](c)
	}else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n",c.Path)
	}
}

