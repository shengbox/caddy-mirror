package mirror

import (
	"bytes"
	"io"
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

// Mirror 是我们的模块结构体
type Mirror struct {
	MirrorTo string `json:"mirror_to,omitempty"` // 镜像目标地址
}

// CaddyModule 返回模块信息
func (Mirror) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.mirror", // 模块的唯一标识
		New: func() caddy.Module { return new(Mirror) },
	}
}

// Provision 设置模块（初始化）
func (d *Mirror) Provision(ctx caddy.Context) error {
	// 这里可以添加初始化逻辑，目前为空
	return nil
}

// ServeHTTP 实现中间件逻辑
func (d Mirror) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	// 异步复制请求到 MirrorTo
	go func(req *http.Request) {
		// 创建一个新的请求副本

		body, err := io.ReadAll(req.Body)
		if err != nil {
			return // 忽略错误，避免影响主请求
		}

		req.Body = io.NopCloser(bytes.NewBuffer(body))

		mirrorReq, err := http.NewRequest(req.Method, "http://"+d.MirrorTo+req.RequestURI, bytes.NewBuffer(body))
		if err != nil {
			return // 忽略错误，避免影响主请求
		}

		// 复制请求头
		for key, values := range req.Header {
			for _, value := range values {
				mirrorReq.Header.Add(key, value)
			}
		}

		// 发送镜像请求
		client := &http.Client{}
		_, _ = client.Do(mirrorReq) // 忽略响应，避免阻塞
	}(r)

	// 主请求继续处理
	return next.ServeHTTP(w, r)
}

// UnmarshalCaddyfile 从 Caddyfile 解析配置
func (d *Mirror) UnmarshalCaddyfile(disp *caddyfile.Dispenser) error {
	for disp.Next() {
		if disp.NextArg() {
			d.MirrorTo = disp.Val()
		}
	}
	return nil
}

// 接口实现
var (
	_ caddyhttp.MiddlewareHandler = (*Mirror)(nil)
	_ caddy.Provisioner           = (*Mirror)(nil)
	_ caddyfile.Unmarshaler       = (*Mirror)(nil)
)

func init() {
	caddy.RegisterModule(Mirror{})

	httpcaddyfile.RegisterHandlerDirective("mirror", func(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
		var d Mirror
		d.UnmarshalCaddyfile(h.Dispenser)
		return &d, nil
	})
}
