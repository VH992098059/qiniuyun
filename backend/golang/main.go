package main

import (
	_ "golang/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"

	"golang/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
