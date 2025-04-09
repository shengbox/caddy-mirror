package main

import (
	_ "github.com/shengbox/caddy-mirror" // 导入我们的模块

	caddycmd "github.com/caddyserver/caddy/v2/cmd"
	_ "github.com/caddyserver/caddy/v2/modules/standard" // 导入标准模块
)

func main() {
	caddycmd.Main()
}
