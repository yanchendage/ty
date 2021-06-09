package main

import (

	"net/http"
	"ty"
)

func main()  {
	r := ty.New()


	r.GET("/fuck", func(c *ty.Context) {
		c.Json(http.StatusOK, ty.J{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	r.GET("/assets/*filepath", func(c *ty.Context) {
		c.Json(http.StatusOK, ty.J{"filepath": c.Param("filepath")})
	})


	v1 := r.Group("/v1")
	v1.GET("/fuck", func(c *ty.Context) {
		v2 := r.Group("/v2")

		c.Json(http.StatusOK, ty.J{"filepath": "nonono"})
	})

	r.Run(":9999")
}
