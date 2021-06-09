package web

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type J map[string]interface{}

type Context struct {
	//res & req
	Writer 		http.ResponseWriter
	Req 		*http.Request
	//req
	Path 		string
	Method 		string
	Params 		map[string]string
	//res
	statusCode 	int
}

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer:   w,
		Req:      r,
		Path:     r.URL.Path,
		Method:   r.Method,
		Params:	  make(map[string]string),
	}
}

func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) Param(key string) string{
	value, _ := c.Params[key]
	return value
}

func (c *Context) SetStatusCode(statusCode int) {
	c.statusCode = statusCode
	c.Writer.WriteHeader(c.statusCode)
}

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) String(statusCode int, format string, value ...interface{}) {
	c.SetHeader("Content-type", "text/plain")
	c.SetStatusCode(statusCode)
	c.Writer.Write([]byte(fmt.Sprintf(format,value...)))
}

func (c * Context) Json(statusCode int, obj interface{})  {
	c.SetHeader("Content-Type", "application/json")
	c.SetStatusCode(statusCode)

	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) Data(statusCode int, data []byte) {
	c.SetStatusCode(statusCode)
	c.Writer.Write(data)
}





