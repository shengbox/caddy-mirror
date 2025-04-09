# Caddy Mirror Module

`caddy-mirror` 是一个自定义的 Caddy 模块，它允许你在将 HTTP 请求反向代理到主目标地址的同时，将请求复制一份发送到另一个指定地址。这个功能适用于需要日志记录、流量审计或将请求复制到备用服务器的场景，且不会影响主请求的处理。

## 功能特点

- 将 HTTP 请求异步镜像到指定地址。
- 镜像过程在后台运行，不干扰主请求的响应。
- 支持智能协议处理（如果未指定协议，默认使用 HTTP）。
- 对镜像地址进行严格验证，确保格式正确。

## 安装

### 1、使用docker镜像
docker pull shengbox/caddy:main


### 2、构建带有 Mirror 模块的 Caddy

你需要通过 `xcaddy` 构建工具将 `caddy-mirror` 模块集成到 Caddy 中。
使用以下 Dockerfile 来构建 Caddy：

```dockerfile
FROM caddy:2.9.1-builder AS builder

RUN xcaddy build \
    --with github.com/shengbox/caddy-mirror@main
```

## 配置
在 Caddyfile 中，你可以使用 mirror 指令配置请求镜像功能。以下是一个示例：
### 示例 Caddyfile


```code
{
    order mirror after respond
}

:80 {
    handle /api/* {
        mirror 192.168.1.188:8088
        reverse_proxy 192.168.1.188:8080
    }
}
```

## 配置说明：
*  { order mirror after respond }：指定 mirror 中间件在 respond 之后执行。  
*  handle /api/*：对 /api/ 路径下的请求应用配置。  
*  mirror 192.168.1.188:8088：将请求镜像到 192.168.1.188:8088。  
*  reverse_proxy 192.168.1.188:8080：将请求反向代理到 192.168.1.188:8080。  

在这个例子中，所有匹配 /api/* 的请求都会被反向代理到 192.168.1.188:8080，同时一份请求副本会被发送到 192.168.1.188:8088。

