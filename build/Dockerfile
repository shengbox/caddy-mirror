# 第一阶段：构建自定义 Caddy
FROM caddy:2.9.1-builder AS builder

RUN xcaddy build \
    --with github.com/shengbox/caddy-mirror@main

# 第二阶段：创建运行时镜像
FROM caddy:2.9.1

# 从 builder 阶段复制构建好的 Caddy 二进制文件
COPY --from=builder /usr/bin/caddy /usr/bin/caddy

CMD ["caddy", "run", "--config", "/etc/caddy/Caddyfile"]