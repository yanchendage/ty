package main

import "github.com/yanchendage/ty/web"

func main()  {
	r := web.New()
	
	r.GET("/", func(c *web.Context) {
		c.Json(200,web.J{"code":200})
	})

	r.Run(":8888")
}
