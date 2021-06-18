package main

import "fmt"

type Length int
var ll Length

func (l Length) fuck()  {
	fmt.Println("fuck")
}

func main()  {
	var l Length
	l = 1
	l.fuck()
	fmt.Println(l)


	//r := web.New()
	//r.Use(logger())
	//
	//r.GET("/", func(c *web.Context) {
	//	c.Json(200,web.J{"code":300})
	//})
	//
	//r.Run(":8888")

	ll = Length(1)



}
