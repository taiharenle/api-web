package middleware

import (
	"api-web/schema"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/handler"
)

func GraphqlHandler() gin.HandlerFunc {
	h := handler.New(&handler.Config{
		Schema:     &schema.Schema,
		Pretty:     true,
		GraphiQL:   true,
		Playground: true,
	})

	// 通过gin封装
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if ip == "::1" {
			ip = "127.0.0.1"
		}
		c.Set("ip", ip)

		h.ContextHandler(c, c.Writer, c.Request)
	}
}
