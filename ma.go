package main

import "github.com/yanchendage/ty/web"

func main()  {
	r := web.New()
	r.Use(logger())

	r.GET("/", func(c *web.Context) {
		c.Json(200,web.J{"code":300})
	})

	r.Run(":8888")
}
