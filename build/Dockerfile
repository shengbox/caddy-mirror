FROM caddy:2.9.1-builder AS builder

RUN xcaddy build \
    --with github.com/shengbox/caddy-mirror@main

FROM caddy:2.9.1

COPY --from=builder /usr/bin/caddy /usr/bin/caddy

CMD ["caddy", "run", "--config", "/etc/caddy/Caddyfile"]