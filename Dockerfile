ARG ALPINE_VERSION=3.15
ARG GO_ALPINE_VERSION=3.15
ARG GO_VERSION=1.17
ARG XCPUTRANSLATE_VERSION=v0.6.0
ARG GOLANGCI_LINT_VERSION=v1.43.0
ARG BUILDPLATFORM=linux/amd64

FROM --platform=${BUILDPLATFORM} qmcgaw/xcputranslate:${XCPUTRANSLATE_VERSION} AS xcputranslate
FROM --platform=${BUILDPLATFORM} qmcgaw/binpot:golangci-lint-${GOLANGCI_LINT_VERSION} AS golangci-lint

FROM --platform=${BUILDPLATFORM} golang:${GO_VERSION}-alpine${GO_ALPINE_VERSION} AS base
COPY --from=xcputranslate /xcputranslate /usr/local/bin/xcputranslate
RUN apk --update add git g++
ENV CGO_ENABLED=0
COPY --from=golangci-lint /bin /go/bin/golangci-lint
WORKDIR /tmp/gobuild
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/ ./cmd/
COPY internal/ ./internal/

FROM --platform=${BUILDPLATFORM} base AS test
# Note on the go race detector:
# - we set CGO_ENABLED=1 to have it enabled
# - we installed g++ to support the race detector
ENV CGO_ENABLED=1
ENTRYPOINT go test -race -coverpkg=./... -coverprofile=coverage.txt -covermode=atomic ./...

FROM --platform=${BUILDPLATFORM} base AS lint
COPY .golangci.yml ./
RUN golangci-lint run --timeout=10m

FROM --platform=${BUILDPLATFORM} base AS tidy
RUN git init && \
    git config user.email ci@localhost && \
    git config user.name ci && \
    git add -A && git commit -m ci && \
    sed -i '/\/\/ indirect/d' go.mod && \
    go mod tidy && \
    git diff --exit-code -- go.mod

FROM --platform=${BUILDPLATFORM} base AS build
ARG TARGETPLATFORM
ARG VERSION=unknown
ARG CREATED="an unknown date"
ARG COMMIT=unknown
RUN GOARCH="$(xcputranslate translate -field arch -targetplatform ${TARGETPLATFORM})" \
    GOARM="$(xcputranslate translate -field arm -targetplatform ${TARGETPLATFORM})" \
    go build -trimpath -ldflags="-s -w \
    -X 'main.version=$VERSION' \
    -X 'main.created=$CREATED' \
    -X 'main.commit=$COMMIT' \
    " -o entrypoint cmd/gluetun/main.go

FROM alpine:${ALPINE_VERSION}
ARG VERSION=unknown
ARG CREATED="an unknown date"
ARG COMMIT=unknown
LABEL \
    org.opencontainers.image.authors="quentin.mcgaw@gmail.com" \
    org.opencontainers.image.created=$CREATED \
    org.opencontainers.image.version=$VERSION \
    org.opencontainers.image.revision=$COMMIT \
    org.opencontainers.image.url="https://github.com/qdm12/gluetun" \
    org.opencontainers.image.documentation="https://github.com/qdm12/gluetun" \
    org.opencontainers.image.source="https://github.com/qdm12/gluetun" \
    org.opencontainers.image.title="VPN swiss-knife like client for multiple VPN providers" \
    org.opencontainers.image.description="VPN swiss-knife like client to tunnel to multiple VPN servers using OpenVPN, IPtables, DNS over TLS, Shadowsocks, an HTTP proxy and Alpine Linux"
ENV VPNSP=pia \
    VPN_TYPE=openvpn \
    # Common VPN options
    VPN_ENDPOINT_IP= \
    VPN_ENDPOINT_PORT= \
    # OpenVPN
    OPENVPN_PROTOCOL=udp \
    OPENVPN_USER= \
    OPENVPN_PASSWORD= \
    OPENVPN_USER_SECRETFILE=/run/secrets/openvpn_user \
    OPENVPN_PASSWORD_SECRETFILE=/run/secrets/openvpn_password \
    OPENVPN_VERSION=2.5 \
    OPENVPN_VERBOSITY=1 \
    OPENVPN_FLAGS= \
    OPENVPN_CIPHER= \
    OPENVPN_AUTH= \
    OPENVPN_PROCESS_USER= \
    OPENVPN_IPV6=off \
    OPENVPN_CUSTOM_CONFIG= \
    OPENVPN_INTERFACE=tun0 \
    # Wireguard
    WIREGUARD_PRIVATE_KEY= \
    WIREGUARD_PRESHARED_KEY= \
    WIREGUARD_PUBLIC_KEY= \
    WIREGUARD_ADDRESS= \
    WIREGUARD_INTERFACE=wg0 \
    # VPN server filtering
    REGION= \
    COUNTRY= \
    CITY= \
    SERVER_HOSTNAME= \
    # # Mullvad only:
    ISP= \
    OWNED_ONLY=no \
    # # Private Internet Access only:
    PIA_ENCRYPTION= \
    PORT_FORWARDING=off \
    PORT_FORWARDING_STATUS_FILE="/tmp/gluetun/forwarded_port" \
    # # Cyberghost only:
    OPENVPN_CLIENTCRT_SECRETFILE=/run/secrets/openvpn_clientcrt \
    OPENVPN_CLIENTKEY_SECRETFILE=/run/secrets/openvpn_clientkey \
    # # Nordvpn only:
    SERVER_NUMBER= \
    # # PIA and ProtonVPN only:
    SERVER_NAME= \
    # # ProtonVPN only:
    FREE_ONLY= \
    # # Surfshark only:
    MULTIHOP_ONLY= \
    # Firewall
    FIREWALL=on \
    FIREWALL_VPN_INPUT_PORTS= \
    FIREWALL_INPUT_PORTS= \
    FIREWALL_OUTBOUND_SUBNETS= \
    FIREWALL_DEBUG=off \
    # Logging
    LOG_LEVEL=info \
    # Health
    HEALTH_SERVER_ADDRESS=127.0.0.1:9999 \
    HEALTH_ADDRESS_TO_PING=github.com \
    HEALTH_VPN_DURATION_INITIAL=6s \
    HEALTH_VPN_DURATION_ADDITION=5s \
    # DNS over TLS
    DOT=on \
    DOT_PROVIDERS=cloudflare \
    DOT_PRIVATE_ADDRESS=127.0.0.1/8,10.0.0.0/8,172.16.0.0/12,192.168.0.0/16,169.254.0.0/16,::1/128,fc00::/7,fe80::/10,::ffff:7f00:1/104,::ffff:a00:0/104,::ffff:a9fe:0/112,::ffff:ac10:0/108,::ffff:c0a8:0/112 \
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
    DNS_PLAINTEXT_ADDRESS=127.0.0.1 \
    DNS_KEEP_NAMESERVER=off \
    # HTTP proxy
    HTTPPROXY= \
    HTTPPROXY_LOG=off \
    HTTPPROXY_LISTENING_ADDRESS=":8888" \
    HTTPPROXY_USER= \
    HTTPPROXY_PASSWORD= \
    HTTPPROXY_USER_SECRETFILE=/run/secrets/httpproxy_user \
    HTTPPROXY_PASSWORD_SECRETFILE=/run/secrets/httpproxy_password \
    # Shadowsocks
    SHADOWSOCKS=off \
    SHADOWSOCKS_LOG=off \
    SHADOWSOCKS_LISTENING_ADDRESS=":8388" \
    SHADOWSOCKS_PASSWORD= \
    SHADOWSOCKS_PASSWORD_SECRETFILE=/run/secrets/shadowsocks_password \
    SHADOWSOCKS_CIPHER=chacha20-ietf-poly1305 \
    # Control server
    HTTP_CONTROL_SERVER_ADDRESS=":8000" \
    # Server data updater
    UPDATER_PERIOD=0 \
    UPDATER_VPN_SERVICE_PROVIDERS= \
    # Public IP
    PUBLICIP_FILE="/tmp/gluetun/ip" \
    PUBLICIP_PERIOD=12h \
    # Pprof
    PPROF_ENABLED=no \
    PPROF_BLOCK_PROFILE_RATE=0 \
    PPROF_MUTEX_PROFILE_RATE=0 \
    PPROF_HTTP_SERVER_ADDRESS=":6060" \
    # Extras
    VERSION_INFORMATION=on \
    TZ= \
    PUID= \
    PGID=
ENTRYPOINT ["/gluetun-entrypoint"]
EXPOSE 8000/tcp 8888/tcp 8388/tcp 8388/udp
HEALTHCHECK --interval=5s --timeout=5s --start-period=10s --retries=1 CMD /gluetun-entrypoint healthcheck
ARG TARGETPLATFORM
RUN apk add --no-cache --update -l apk-tools && \
    apk add --no-cache --update -X "https://dl-cdn.alpinelinux.org/alpine/v3.12/main" openvpn==2.4.11-r0 && \
    mv /usr/sbin/openvpn /usr/sbin/openvpn2.4 && \
    apk del openvpn && \
    apk add --no-cache --update openvpn ca-certificates iptables ip6tables unbound tzdata && \
    mv /usr/sbin/openvpn /usr/sbin/openvpn2.5 && \
    # Fix vulnerability issue
    apk add --no-cache --update busybox && \
    rm -rf /var/cache/apk/* /etc/unbound/* /usr/sbin/unbound-* /etc/openvpn/*.sh /usr/lib/openvpn/plugins/openvpn-plugin-down-root.so && \
    deluser openvpn && \
    deluser unbound && \
    mkdir /gluetun
COPY --from=build /tmp/gobuild/entrypoint /gluetun-entrypoint
