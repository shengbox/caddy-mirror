package mirror

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

type Mirror struct {
	MirrorTo  string `json:"mirror_to,omitempty"`
	mirrorURL *url.URL
}

func (Mirror) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.mirror",
		New: func() caddy.Module { return new(Mirror) },
	}
}

func (d *Mirror) Provision(ctx caddy.Context) error {
	input := strings.TrimSpace(d.MirrorTo)
	if input == "" {
		return errors.New("mirror_to must not be empty")
	}
	// 智能协议处理
	if !strings.Contains(input, "://") {
		input = "http://" + input // 默认使用HTTP协议
	}
	// 严格验证URL合法性
	u, err := url.Parse(input)
	if err != nil {
		return fmt.Errorf("invalid mirror_to URL: %w", err)
	}
	d.mirrorURL = u
	return nil
}

func (m Mirror) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {

	go func(req *http.Request) {
		// 创建一个新的请求副本
		body, err := io.ReadAll(req.Body)
		if err != nil {
			return
		}

		dest := m.mirrorURL.ResolveReference(req.URL).String()
		req.Body = io.NopCloser(bytes.NewBuffer(body))
		mirrorReq, err := http.NewRequest(req.Method, dest, bytes.NewBuffer(body))
		if err != nil {
			return // 忽略错误，避免影响主请求
		}

		// 复制请求头
		for key, values := range req.Header {
			for _, value := range values {
				mirrorReq.Header.Add(key, value)
			}
		}

		client := &http.Client{}
		_, _ = client.Do(mirrorReq)
	}(r)

	return next.ServeHTTP(w, r)
}

func (d *Mirror) UnmarshalCaddyfile(disp *caddyfile.Dispenser) error {
	for disp.Next() {
		if disp.NextArg() {
			d.MirrorTo = disp.Val()
		}
	}
	return nil
}

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
