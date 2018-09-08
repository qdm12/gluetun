FROM alpine:3.8
LABEL maintainer="quentin.mcgaw@gmail.com" \
      description="VPN client to private internet access servers using OpenVPN, Alpine, IPtables firewall and Cloudflare 1.1.1.1 DNS over TLS" \
      download="5.7MB" \
      size="13.5MB" \
      ram="12MB" \
      cpu_usage="Low" \
      github="https://github.com/qdm12/private-internet-access-docker"
RUN apk add -q --progress --no-cache --update openvpn unbound ca-certificates iptables && \
    apk add -q --progress --no-cache --update --virtual=build-dependencies unzip && \
    mkdir /openvpn-udp-normal /openvpn-udp-strong /openvpn-tcp-normal /openvpn-tcp-strong && \
    wget -q https://www.privateinternetaccess.com/openvpn/openvpn.zip \
            https://www.privateinternetaccess.com/openvpn/openvpn-strong.zip \
            https://www.privateinternetaccess.com/openvpn/openvpn-tcp.zip \
            https://www.privateinternetaccess.com/openvpn/openvpn-strong-tcp.zip && \
    unzip -q openvpn.zip -d /openvpn-udp-normal && \
    unzip -q openvpn-strong.zip -d /openvpn-udp-strong && \
    unzip -q openvpn-tcp.zip -d /openvpn-tcp-normal && \
    unzip -q openvpn-strong-tcp.zip -d /openvpn-tcp-strong && \
    apk del -q --progress --purge build-dependencies && \
    rm -rf /*.zip /etc/unbound/unbound.conf /var/cache/apk/*
COPY unbound.conf /etc/unbound/unbound.conf
HEALTHCHECK --interval=10m --timeout=10s --start-period=10s --retries=1 \
            CMD export OLD_VPN_IP="$NEW_VPN_IP" && \
                export NEW_VPN_IP=$(wget -qqO- 'https://duckduckgo.com/?q=what+is+my+ip' | grep -ow 'Your IP address is [0-9.]*[0-9]' | grep -ow '[0-9][0-9.]*') && \
                [ "$NEW_VPN_IP" != "$INITIAL_IP" ] && [ "$NEW_VPN_IP" != "$OLD_VPN_IP" ] || exit 1
ENV ENCRYPTION=strong \
    PROTOCOL=tcp \
    REGION=Germany
COPY entrypoint.sh /
RUN chmod +x /entrypoint.sh
ENTRYPOINT /entrypoint.sh