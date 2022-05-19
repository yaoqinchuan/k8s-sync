package main

import (
	"context"
	"fmt"
	_ "mygogf/internal/packed"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gctx"
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
	Port   int    `v:"required" short:"p" name:"port" arg:"true" brief:"port of http server"`
	Mode   string `v:"required" short:"m" name:"mode" arg:"true" brief:"run mode: dev test prod"`
}
type cCmdOutput struct{}

var args inputArgs = inputArgs{}

func (c *cMain) CmdInit(ctx context.Context, in cCmdInput) (out *cCmdOutput, err error) {
	args.Port = in.Port
	args.Mode = in.Mode
	s := g.Server("sync-k8s")
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Write(fmt.Printf("app %v run with mode %v on port %v", args.Mode, args.Port))
	})
	s.SetPort(args.Port)
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
