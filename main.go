package main

import (
	"context"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gctx"
	"k8s-sync/internal/controller/v1"
	middleware "k8s-sync/internal/service/middlewares/before"
	_ "k8s-sync/internal/utils"
)

type inputArgs struct {
	Port int
	Mode string
}

type cMain struct {
	g.Meta `name:"run" brief:"app option"`
}

type cCmdInput struct {
	g.Meta `name:"app" brief:"app option"`
	Port   int    `v:"required" short:"p" name:"port" arg:"true" brief:"port of http server,default 8080"`
	Mode   string `v:"required" short:"m" name:"mode" arg:"true" brief:"run mode: dev test prod, default dev"`
}
type cCmdOutput struct{}

var args = inputArgs{
	Port: 8080,
	Mode: "dev",
}

func corsMiddleware(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}

func apiHandlerRegister(group *ghttp.RouterGroup) {
	group.Middleware(middleware.AuthMiddleware)
	group.Group("/admin", func(group *ghttp.RouterGroup) {
		controller.AccountApiHandlerRegister(group)
	})
	group.Group("/maintenance", func(group *ghttp.RouterGroup) {
		controller.MaintenanceApiHandlerRegister(group)
	})
}
func heathCheckRegister(group *ghttp.RouterGroup) {
	group.GET("healthy", func(r *ghttp.Request) {
		r.Response.Write("ok")
	})
}
func (c *cMain) CmdInit(ctx context.Context, in cCmdInput) (out *cCmdOutput, err error) {
	args.Port = in.Port
	args.Mode = in.Mode
	s := g.Server("sync-k8s")
	s.SetPort(args.Port)
	s.Group("/api/v1", func(group *ghttp.RouterGroup) {
		group.Middleware(middleware.ApiRequestRecord)
		group.Middleware(corsMiddleware)
		apiHandlerRegister(group)
		heathCheckRegister(group)
	})

	s.Run()
	return
}
func main() {
	cmd, err := gcmd.NewFromObject(cMain{})
	if err != nil {
		panic(err)
	}
	cmd.Run(gctx.New())
}
