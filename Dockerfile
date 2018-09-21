FROM alpine:3.8
LABEL maintainer="quentin.mcgaw@gmail.com" \
      description="VPN client to private internet access servers using OpenVPN, IPtables firewall, DNS over TLS with Unbound and Alpine Linux" \
      download="???MB" \
      size="15.7MB" \
      ram="13MB" \
      cpu_usage="Low" \
      github="https://github.com/qdm12/private-internet-access-docker"
HEALTHCHECK --interval=1m --timeout=10s --start-period=10s --retries=1 \
            CMD export OLD_VPN_IP="$NEW_VPN_IP" && \
                export NEW_VPN_IP=$(wget -qqO- 'https://duckduckgo.com/?q=what+is+my+ip' | grep -ow 'Your IP address is [0-9.]*[0-9]' | grep -ow '[0-9][0-9.]*') && \
                [ "$NEW_VPN_IP" != "$INITIAL_IP" ] && [ "$NEW_VPN_IP" != "$OLD_VPN_IP" ] || exit 1
ENV ENCRYPTION=strong \
    PROTOCOL=tcp \
    REGION=Germany
RUN apk add -q --progress --no-cache --update openvpn ca-certificates iptables ip6tables unbound && \
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
    rm -rf /*.zip /var/cache/apk/* /etc/unbound/unbound.conf && \
    addgroup -S nonrootusers && adduser -S nonrootuser -G nonrootusers
COPY unbound.conf /etc/unbound/unbound.conf
COPY entrypoint.sh /
ENTRYPOINT /entrypoint.sh