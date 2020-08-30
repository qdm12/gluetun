ARG ALPINE_VERSION=3.12
ARG GO_VERSION=1.15

FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS builder
RUN apk --update add git
ENV CGO_ENABLED=0
ARG GOLANGCI_LINT_VERSION=v1.30.0
RUN wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s ${GOLANGCI_LINT_VERSION}
WORKDIR /tmp/gobuild
COPY .golangci.yml .
COPY go.mod go.sum ./
RUN go mod download 2>&1
ARG VERSION=unknown
ARG BUILD_DATE="an unknown date"
ARG COMMIT=unknown
COPY cmd/gluetun/main.go .
COPY internal/ ./internal/
RUN go test ./...
RUN golangci-lint run --timeout=10m
RUN go build -trimpath -ldflags="-s -w \
    -X 'main.version=$VERSION' \
    -X 'main.buildDate=$BUILD_DATE' \
    -X 'main.commit=$COMMIT' \
    " -o entrypoint main.go

FROM alpine:${ALPINE_VERSION}
ARG VERSION=unknown
ARG BUILD_DATE="an unknown date"
ARG COMMIT=unknown
LABEL \
    org.opencontainers.image.authors="quentin.mcgaw@gmail.com" \
    org.opencontainers.image.created=$BUILD_DATE \
    org.opencontainers.image.version=$VERSION \
    org.opencontainers.image.revision=$COMMIT \
    org.opencontainers.image.url="https://github.com/qdm12/gluetun" \
    org.opencontainers.image.documentation="https://github.com/qdm12/gluetun" \
    org.opencontainers.image.source="https://github.com/qdm12/gluetun" \
    org.opencontainers.image.title="VPN client for PIA, Mullvad, Windscribe, Surfshark and Cyberghost" \
    org.opencontainers.image.description="VPN client to tunnel to PIA, Mullvad, Windscribe, Surfshark and Cyberghost servers using OpenVPN, IPtables, DNS over TLS and Alpine Linux"
ENV VPNSP=pia \
    VERSION_INFORMATION=on \
    PROTOCOL=udp \
    OPENVPN_VERBOSITY=1 \
    OPENVPN_ROOT=no \
    OPENVPN_TARGET_IP= \
    TZ= \
    UID=1000 \
    GID=1000 \
    IP_STATUS_FILE="/tmp/gluetun/ip" \
    # PIA, Windscribe, Surfshark, Cyberghost, Vyprvpn, NordVPN, PureVPN only
    USER= \
    PASSWORD= \
    REGION= \
    # PIA only
    PIA_ENCRYPTION=strong \
    PORT_FORWARDING=off \
    PORT_FORWARDING_STATUS_FILE="/tmp/gluetun/forwarded_port" \
    # Mullvad and PureVPN only
    COUNTRY= \
    CITY= \
    # Mullvad only
    ISP= \
    # Mullvad and Windscribe only
    PORT= \
    # Cyberghost only
    CYBERGHOST_GROUP="Premium UDP Europe" \
    # NordVPN only
    SERVER_NUMBER= \
    # Openvpn
    OPENVPN_CIPHER= \
    OPENVPN_AUTH= \
    # DNS over TLS
    DOT=on \
    DOT_PROVIDERS=cloudflare \
    DOT_PRIVATE_ADDRESS=127.0.0.1/8,10.0.0.0/8,172.16.0.0/12,192.168.0.0/16,169.254.0.0/16,::1/128,fc00::/7,fe80::/10,::ffff:0:0/96 \
    DOT_VERBOSITY=1 \
    DOT_VERBOSITY_DETAILS=0 \
    DOT_VALIDATION_LOGLEVEL=0 \
    DOT_CACHING=on \
    DOT_IPV6=off \
    BLOCK_MALICIOUS=on \
    BLOCK_SURVEILLANCE=off \
    BLOCK_ADS=off \
    UNBLOCK= \
    DNS_UPDATE_PERIOD=24h \
    DNS_PLAINTEXT_ADDRESS=1.1.1.1 \
    DNS_KEEP_NAMESERVER=off \
    # Firewall
    FIREWALL=on \
    EXTRA_SUBNETS= \
    FIREWALL_VPN_INPUT_PORTS= \
    FIREWALL_DEBUG=off \
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
    SHADOWSOCKS_PASSWORD= \
    SHADOWSOCKS_METHOD=chacha20-ietf-poly1305
ENTRYPOINT ["/entrypoint"]
EXPOSE 8000/tcp 8888/tcp 8388/tcp 8388/udp
HEALTHCHECK --interval=10m --timeout=10s --start-period=30s --retries=2 CMD /entrypoint healthcheck
RUN apk add -q --progress --no-cache --update openvpn ca-certificates iptables ip6tables unbound tinyproxy tzdata && \
    rm -rf /var/cache/apk/* /etc/unbound/* /usr/sbin/unbound-* /etc/tinyproxy/tinyproxy.conf && \
    deluser openvpn && \
    deluser tinyproxy && \
    deluser unbound && \
    mkdir /gluetun
COPY --from=builder /tmp/gobuild/entrypoint /entrypoint
