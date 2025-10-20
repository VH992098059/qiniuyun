package cmd

import (
	"context"
	"golang/internal/controller/voice"
	"golang/internal/controller/ws"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()
			//是否允许跨域操作
			s.Use(func(r *ghttp.Request) {
				r.Response.CORSDefault()
				r.Middleware.Next()
			})
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.GET("/ws", func(r *ghttp.Request) {
					ws.WsHandler(r.Response.ResponseWriter, r.Request)
				})
			})
			s.Group("/client", func(group *ghttp.RouterGroup) {
				group.Middleware(ghttp.MiddlewareHandlerResponse)
				group.Bind(
					voice.NewV1(),
				)
			})
			s.Run()
			return nil
		},
	}
)
