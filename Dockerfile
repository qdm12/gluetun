FROM alpine:3.7
LABEL maintainer="quentin.mcgaw@gmail.com" \
      description="VPN client to private internet access servers using OpenVPN, Alpine and Cloudflare 1.1.1.1 DNS over TLS" \
      download="5.6MB" \
      size="13.5MB" \
      ram="12MB" \
      cpu_usage="Low" \
      github="https://github.com/qdm12/private-internet-access-docker"
RUN apk add -q --progress --no-cache --update openvpn unbound ca-certificates && \
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
HEALTHCHECK --interval=10m --timeout=5s --start-period=10s --retries=1 \
            CMD VPNCITY=$(wget -qO- -T 2 https://ipinfo.io/city); \
                VPNORGANIZATION=$(wget -qO- -T 2 https://ipinfo.io/org); \
            printf "\nCity: $VPNCITY\nOrganization: $VPNORGANIZATION"; \
            [ "$VPNCITY" != "$CITY" ] || [ "$VPNORGANIZATION" != "$ORGANIZATION" ] || exit 1
ENV ENCRYPTION=strong \
    PROTOCOL=tcp \
    REGION=Switzerland
COPY entrypoint.sh /
RUN chmod +x /entrypoint.sh
ENTRYPOINT /entrypoint.sh
