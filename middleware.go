package main

import (
	"github.com/yanchendage/ty/web"
	"log"
	"time"
)

func logger() web.HandlerFunc {
	return func(c *web.Context) {
		t := time.Now()
		c.Next()
		log.Println(123)

		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}
