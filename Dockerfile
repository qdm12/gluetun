ARG ALPINE_VERSION=3.11
ARG GO_VERSION=1.14

FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS builder
RUN apk --update add git
WORKDIR /tmp/gobuild
ENV CGO_ENABLED=0
COPY go.mod go.sum ./
RUN go mod download 2>&1
COPY cmd/main.go .
COPY internal/ ./internal/
RUN go test ./...
RUN go build -ldflags="-s -w" -o entrypoint main.go

FROM alpine:${ALPINE_VERSION}
ARG VERSION
ARG BUILD_DATE
ARG VCS_REF
ENV VERSION=$VERSION \
    BUILD_DATE=$BUILD_DATE \
    VCS_REF=$VCS_REF
LABEL \
    org.opencontainers.image.authors="quentin.mcgaw@gmail.com" \
    org.opencontainers.image.created=$BUILD_DATE \
    org.opencontainers.image.version=$VERSION \
    org.opencontainers.image.revision=$VCS_REF \
    org.opencontainers.image.url="https://github.com/qdm12/private-internet-access-docker" \
    org.opencontainers.image.documentation="https://github.com/qdm12/private-internet-access-docker" \
    org.opencontainers.image.source="https://github.com/qdm12/private-internet-access-docker" \
    org.opencontainers.image.title="PIA client" \
    org.opencontainers.image.description="VPN client to tunnel to private internet access servers using OpenVPN, IPtables, DNS over TLS and Alpine Linux"
ENV VPNSP=pia \
    USER= \
    PROTOCOL=udp \
    OPENVPN_VERBOSITY=1 \
    OPENVPN_ROOT=no \
    OPENVPN_TARGET_IP= \
    TZ= \
    # PIA only
    PASSWORD= \
    REGION="CA Montreal" \
    PIA_ENCRYPTION=strong \
    PORT_FORWARDING=off \
    PORT_FORWARDING_STATUS_FILE="/forwarded_port" \
    # Mullvad only
    COUNTRY=Sweden \
    CITY= \
    ISP= \
    PORT= \
    # DNS over TLS
    DOT=on \
    DOT_PROVIDERS=cloudflare \
    DOT_PRIVATE_ADDRESS=127.0.0.1/8,10.0.0.0/8,172.16.0.0/12,192.168.0.0/16,169.254.0.0/16,::1/128,fc00::/7,fe80::/10,::ffff:0:0/96 \
    DOT_VERBOSITY=1 \
    DOT_VERBOSITY_DETAILS=0 \
    DOT_VALIDATION_LOGLEVEL=0 \
    DOT_CACHING=on \
    DOT_IPV6=on \
    BLOCK_MALICIOUS=on \
    BLOCK_SURVEILLANCE=off \
    BLOCK_ADS=off \
    UNBLOCK= \
    # Firewall
    EXTRA_SUBNETS= \
    # Tinyproxy
    TINYPROXY=off \
    TINYPROXY_LOG=Info \
    TINYPROXY_PORT=8888 \
    TINYPROXY_USER= \
    TINYPROXY_PASSWORD= \
    # Shadowsocks
    SHADOWSOCKS=off \
    SHADOWSOCKS_LOG=off \
    SHADOWSOCKS_PORT=8388 \
    SHADOWSOCKS_PASSWORD=
ENTRYPOINT /entrypoint
EXPOSE 8888/tcp 8388/tcp 8388/udp
HEALTHCHECK --interval=3m --timeout=3s --start-period=20s --retries=1 CMD /entrypoint healthcheck
RUN apk add -q --progress --no-cache --update openvpn ca-certificates iptables unbound tinyproxy tzdata && \
    echo "http://dl-cdn.alpinelinux.org/alpine/edge/testing" >> /etc/apk/repositories && \
    apk add -q --progress --no-cache --update shadowsocks-libev && \
    rm -rf /*.zip /var/cache/apk/* /etc/unbound/* /usr/sbin/unbound-anchor /usr/sbin/unbound-checkconf /usr/sbin/unbound-control /usr/sbin/unbound-control-setup /usr/sbin/unbound-host /etc/tinyproxy/tinyproxy.conf && \
    adduser nonrootuser -D -H --uid 1000 && \
    chown nonrootuser -R /etc/unbound /etc/tinyproxy && \
    chmod 700 /etc/unbound /etc/tinyproxy
COPY --from=builder --chown=1000:1000 /tmp/gobuild/entrypoint /entrypoint