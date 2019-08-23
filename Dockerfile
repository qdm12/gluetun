# Latest version of ubuntu
FROM ubuntu:18.04

ENV USER= \
    PASSWORD= \
    ENCRYPTION=strong \
    PROTOCOL=udp \
    REGION="UK London" \
    NONROOT=no \
    EXTRA_SUBNETS= \
    PORT_FORWARDING=false \
    PROXY=on \
    PROXY_LOG_LEVEL=Critical \
    PROXY_PORT=8888 \
    PROXY_USER= \
    PROXY_PASSWORD=

# Start point for docker
ENTRYPOINT /entrypoint.sh

# Port for qBittorrent
EXPOSE 8888

HEALTHCHECK --interval=3m --timeout=3s --start-period=20s --retries=1 CMD /healthcheck.sh

# Ok lets install everything
RUN apt-get update \
    && set -x \
    && apt-get install -qq --no-install-recommends -y qbittorrent openvpn wget ca-certificates iptables unzip \
    wget -q https://www.privateinternetaccess.com/openvpn/openvpn.zip \
    https://www.privateinternetaccess.com/openvpn/openvpn-strong.zip \
    https://www.privateinternetaccess.com/openvpn/openvpn-tcp.zip \
    https://www.privateinternetaccess.com/openvpn/openvpn-strong-tcp.zip && \
    mkdir -p /openvpn/target && \
    unzip -q openvpn.zip -d /openvpn/udp-normal && \
    unzip -q openvpn-strong.zip -d /openvpn/udp-strong && \
    unzip -q openvpn-tcp.zip -d /openvpn/tcp-normal && \
    unzip -q openvpn-strong-tcp.zip -d /openvpn/tcp-strong && \
    && apt-get purge -y -qq unzip \
    && apt-get clean -qq
    rm -rf /*.zip /etc/tinyproxy/tinyproxy.conf && \

COPY tinyproxy.conf /etc/tinyproxy/tinyproxy.conf
COPY entrypoint.sh healthcheck.sh portforward.sh /
RUN chown nonrootuser -R /etc/tinyproxy && \
    chmod 700 /etc/tinyproxy && \
    chmod 600 /etc/tinyproxy/tinyproxy.conf && \
    chmod 500 /entrypoint.sh /healthcheck.sh /portforward.sh
