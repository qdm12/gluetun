ARG ALPINE_VERSION=3.12
ARG GO_VERSION=1.15
ARG BUILDPLATFORM=linux/amd64

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS base
RUN apk --update add git
ENV CGO_ENABLED=0
WORKDIR /tmp/gobuild
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/ ./cmd/
COPY internal/ ./internal/

FROM --platform=$BUILDPLATFORM base AS test
# Note on the go race detector:
# - we set CGO_ENABLED=1 to have it enabled
# - we install g++ to support the race detector
ENV CGO_ENABLED=1
RUN apk --update --no-cache add g++

FROM --platform=$BUILDPLATFORM base AS lint
ARG GOLANGCI_LINT_VERSION=v1.35.2
RUN wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
    sh -s -- -b /usr/local/bin ${GOLANGCI_LINT_VERSION}
COPY .golangci.yml ./
RUN golangci-lint run --timeout=10m

FROM --platform=$BUILDPLATFORM base AS tidy
RUN git init && \
    git config user.email ci@localhost && \
    git config user.name ci && \
    git add -A && git commit -m ci && \
    sed -i '/\/\/ indirect/d' go.mod && \
    go mod tidy && \
    git diff --exit-code -- go.mod

FROM --platform=$BUILDPLATFORM base AS build
COPY --from=qmcgaw/xcputranslate:v0.4.0 /xcputranslate /usr/local/bin/xcputranslate
ARG TARGETPLATFORM
ARG VERSION=unknown
ARG BUILD_DATE="an unknown date"
ARG COMMIT=unknown
RUN GOARCH="$(xcputranslate -field arch -targetplatform ${TARGETPLATFORM})" \
    GOARM="$(xcputranslate -field arm -targetplatform ${TARGETPLATFORM})" \
    go build -trimpath -ldflags="-s -w \
    -X 'main.version=$VERSION' \
    -X 'main.buildDate=$BUILD_DATE' \
    -X 'main.commit=$COMMIT' \
    " -o entrypoint cmd/gluetun/main.go

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
    org.opencontainers.image.title="VPN swiss-knife like client for multiple VPN providers" \
    org.opencontainers.image.description="VPN swiss-knife like client to tunnel to multiple VPN servers using OpenVPN, IPtables, DNS over TLS, Shadowsocks, an HTTP proxy and Alpine Linux"
ENV VPNSP=pia \
    VERSION_INFORMATION=on \
    PROTOCOL=udp \
    OPENVPN_VERBOSITY=1 \
    OPENVPN_ROOT=no \
    OPENVPN_TARGET_IP= \
    OPENVPN_IPV6=off \
    TZ= \
    PUID= \
    PGID= \
    PUBLICIP_FILE="/tmp/gluetun/ip" \
    # PIA, Windscribe, Surfshark, Cyberghost, Vyprvpn, NordVPN, PureVPN only
    OPENVPN_USER= \
    OPENVPN_PASSWORD= \
    USER_SECRETFILE=/run/secrets/openvpn_user \
    PASSWORD_SECRETFILE=/run/secrets/openvpn_password \
    REGION= \
    # PIA only
    PIA_ENCRYPTION=strong \
    PORT_FORWARDING=off \
    PORT_FORWARDING_STATUS_FILE="/tmp/gluetun/forwarded_port" \
    # Mullvad and PureVPN only
    COUNTRY= \
    # Mullvad, PureVPN, Windscribe only
    CITY= \
    # Windscribe only
    SERVER_HOSTNAME= \
    # Mullvad only
    ISP= \
    OWNED=no \
    # Mullvad and Windscribe only
    PORT= \
    # Cyberghost only
    CYBERGHOST_GROUP="Premium UDP Europe" \
    OPENVPN_CLIENTCRT_SECRETFILE=/run/secrets/openvpn_clientcrt \
    OPENVPN_CLIENTKEY_SECRETFILE=/run/secrets/openvpn_clientkey \
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
    FIREWALL_VPN_INPUT_PORTS= \
    FIREWALL_INPUT_PORTS= \
    FIREWALL_OUTBOUND_SUBNETS= \
    FIREWALL_DEBUG=off \
    # HTTP proxy
    HTTPPROXY= \
    HTTPPROXY_LOG=off \
    HTTPPROXY_PORT=8888 \
    HTTPPROXY_USER= \
    HTTPPROXY_PASSWORD= \
    HTTPPROXY_USER_SECRETFILE=/run/secrets/httpproxy_user \
    HTTPPROXY_PASSWORD_SECRETFILE=/run/secrets/httpproxy_password \
    # Shadowsocks
    SHADOWSOCKS=off \
    SHADOWSOCKS_LOG=off \
    SHADOWSOCKS_PORT=8388 \
    SHADOWSOCKS_PASSWORD= \
    SHADOWSOCKS_PASSWORD_SECRETFILE=/run/secrets/shadowsocks_password \
    SHADOWSOCKS_METHOD=chacha20-ietf-poly1305 \
    UPDATER_PERIOD=0
ENTRYPOINT ["/entrypoint"]
EXPOSE 8000/tcp 8888/tcp 8388/tcp 8388/udp
HEALTHCHECK --interval=5s --timeout=5s --start-period=10s --retries=1 CMD /entrypoint healthcheck
RUN apk add -q --progress --no-cache --update openvpn ca-certificates iptables ip6tables unbound tzdata && \
    rm -rf /var/cache/apk/* /etc/unbound/* /usr/sbin/unbound-* && \
    deluser openvpn && \
    deluser unbound && \
    mkdir /gluetun
# TODO remove once SAN is added to PIA servers certificates, see https://github.com/pia-foss/manual-connections/issues/10
ENV GODEBUG=x509ignoreCN=0
COPY --from=build /tmp/gobuild/entrypoint /entrypoint
