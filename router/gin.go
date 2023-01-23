package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pallat/todoapi/todo"
)

type MyContext struct {
	*gin.Context
}

func NewMyContext(c *gin.Context) *MyContext {
	return &MyContext{Context: c}
}

func (c *MyContext) Bind(v interface{}) error {
	return c.Context.ShouldBindJSON(v)
}

func (c *MyContext) JSON(statuscode int, obj interface{}) {
	c.Context.JSON(statuscode, obj)
}

func (c *MyContext) TransactionID() string {
	return c.Request.Header.Get("TransactionID")
}

func (c *MyContext) Audience() string {
	if aud, ok := c.Get("aud"); ok {
		if s, ok := aud.(string); ok {
			return s
		}
	}
	return ""
}

func (c *MyContext) Param(v string) string {
	return c.Context.Param(v)
}

func NewGinHandler(handler func(todo.Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler(NewMyContext(c))
	}
}

type MyRouter struct {
	*gin.Engine
}

func NewMyRouter() *MyRouter {
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{
		"http://localhost:8080",
	}
	config.AllowHeaders = []string{
		"Origin",
		"Authorization",
		"TansactionID",
	}
	r.Use(cors.New(config))

	return &MyRouter{Engine: r}
}

func (r *MyRouter) GET(path string, handler func(todo.Context)) {
	r.Engine.GET(path, NewGinHandler(handler))
}

func (r *MyRouter) POST(path string, handler func(todo.Context)) {
	r.Engine.POST(path, NewGinHandler(handler))
}

func (r *MyRouter) PUT(path string, handler func(todo.Context)) {
	r.Engine.PUT(path, NewGinHandler(handler))
}

func (r *MyRouter) DELETE(path string, handler func(todo.Context)) {
	r.Engine.DELETE(path, NewGinHandler(handler))
}
