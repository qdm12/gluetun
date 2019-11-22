ARG ALPINE_VERSION=3.10

FROM alpine:${ALPINE_VERSION}
ARG BUILD_DATE
ARG VCS_REF
LABEL \
    org.opencontainers.image.authors="quentin.mcgaw@gmail.com" \
    org.opencontainers.image.created=$BUILD_DATE \
    org.opencontainers.image.version="" \
    org.opencontainers.image.revision=$VCS_REF \
    org.opencontainers.image.url="https://github.com/qdm12/private-internet-access-docker" \
    org.opencontainers.image.documentation="https://github.com/qdm12/private-internet-access-docker" \
    org.opencontainers.image.source="https://github.com/qdm12/private-internet-access-docker" \
    org.opencontainers.image.title="PIA client" \
    org.opencontainers.image.description="VPN client to tunnel to private internet access servers using OpenVPN, IPtables, DNS over TLS and Alpine Linux" \
    image-size="23.3MB" \
    ram-usage="13MB to 80MB" \
    cpu-usage="Low to Medium"
ENV USER= \
    PASSWORD= \
    ENCRYPTION=strong \
    PROTOCOL=udp \
    REGION="CA Montreal" \
    NONROOT=no \
    DOT=on \
    BLOCK_MALICIOUS=off \
    BLOCK_NSA=off \
    UNBLOCK= \
    EXTRA_SUBNETS= \
    PORT_FORWARDING=off \
    PORT_FORWARDING_STATUS_FILE="/forwarded_port" \
    TINYPROXY=off \
    TINYPROXY_LOG=Critical \
    TINYPROXY_PORT=8888 \
    TINYPROXY_USER= \
    TINYPROXY_PASSWORD= \
    SHADOWSOCKS=off \
    SHADOWSOCKS_LOG=on \
    SHADOWSOCKS_PORT=8388 \
    SHADOWSOCKS_PASSWORD= \
    TZ=
ENTRYPOINT /entrypoint.sh
EXPOSE 8888/tcp 8388/tcp 8388/udp
HEALTHCHECK --interval=3m --timeout=3s --start-period=20s --retries=1 CMD /healthcheck.sh
RUN apk add -q --progress --no-cache --update openvpn wget ca-certificates iptables unbound unzip tinyproxy jq tzdata && \
    echo "http://dl-cdn.alpinelinux.org/alpine/edge/testing" >> /etc/apk/repositories && \
    apk add -q --progress --no-cache --update shadowsocks-libev && \
    wget -q https://www.privateinternetaccess.com/openvpn/openvpn.zip \
    https://www.privateinternetaccess.com/openvpn/openvpn-strong.zip \
    https://www.privateinternetaccess.com/openvpn/openvpn-tcp.zip \
    https://www.privateinternetaccess.com/openvpn/openvpn-strong-tcp.zip && \
    mkdir -p /openvpn/target && \
    unzip -q openvpn.zip -d /openvpn/udp-normal && \
    unzip -q openvpn-strong.zip -d /openvpn/udp-strong && \
    unzip -q openvpn-tcp.zip -d /openvpn/tcp-normal && \
    unzip -q openvpn-strong-tcp.zip -d /openvpn/tcp-strong && \
    apk del -q --progress --purge unzip && \
    rm -rf /*.zip /var/cache/apk/* /etc/unbound/* /usr/sbin/unbound-anchor /usr/sbin/unbound-checkconf /usr/sbin/unbound-control /usr/sbin/unbound-control-setup /usr/sbin/unbound-host /etc/tinyproxy/tinyproxy.conf && \
    adduser nonrootuser -D -H --uid 1000 && \
    wget -q https://raw.githubusercontent.com/qdm12/files/master/named.root.updated -O /etc/unbound/root.hints && \
    wget -q https://raw.githubusercontent.com/qdm12/files/master/root.key.updated -O /etc/unbound/root.key && \
    cd /tmp && \
    wget -q https://raw.githubusercontent.com/qdm12/files/master/malicious-hostnames.updated -O malicious-hostnames && \
    wget -q https://raw.githubusercontent.com/qdm12/files/master/surveillance-hostnames.updated -O nsa-hostnames && \
    wget -q https://raw.githubusercontent.com/qdm12/files/master/malicious-ips.updated -O malicious-ips && \
    while read hostname; do echo "local-zone: \""$hostname"\" static" >> blocks-malicious.conf; done < malicious-hostnames && \
    while read ip; do echo "private-address: $ip" >> blocks-malicious.conf; done < malicious-ips && \
    tar -cjf /etc/unbound/blocks-malicious.bz2 blocks-malicious.conf && \
    while read hostname; do echo "local-zone: \""$hostname"\" static" >> blocks-nsa.conf; done < nsa-hostnames && \
    tar -cjf /etc/unbound/blocks-nsa.bz2 blocks-nsa.conf && \
    rm -f /tmp/*
COPY unbound.conf /etc/unbound/unbound.conf
COPY tinyproxy.conf /etc/tinyproxy/tinyproxy.conf
COPY shadowsocks.json /etc/shadowsocks.json
COPY entrypoint.sh healthcheck.sh portforward.sh /
RUN chown nonrootuser -R /etc/unbound /etc/tinyproxy && \
    chmod 700 /etc/unbound /etc/tinyproxy && \
    chmod 600 /etc/unbound/unbound.conf /etc/tinyproxy/tinyproxy.conf /etc/shadowsocks.json && \
    chmod 500 /entrypoint.sh /healthcheck.sh /portforward.sh && \
    chmod 400 /etc/unbound/root.hints /etc/unbound/root.key /etc/unbound/*.bz2
