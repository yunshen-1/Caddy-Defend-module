FROM caddy:builder AS builder

COPY . /defender

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    xcaddy build \
    --with pkg.jsn.cam/caddy-defender=/defender

FROM caddy:latest
COPY --from=builder /usr/bin/caddy /usr/bin/caddy
