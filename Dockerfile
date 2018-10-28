FROM alpine:3.8
LABEL maintainer="quentin.mcgaw@gmail.com" \
      description="VPN client to private internet access servers using OpenVPN, IPtables firewall, DNS over TLS with Unbound and Alpine Linux" \
      download="6.6MB" \
      size="15.7MB" \
      ram="13MB" \
      cpu_usage="Low" \
      github="https://github.com/qdm12/private-internet-access-docker"
ENV ENCRYPTION=strong \
    PROTOCOL=tcp \
    REGION="CA Montreal" \
    BLOCK_MALICIOUS=off
HEALTHCHECK --interval=5m --timeout=15s --start-period=10s --retries=2 \
            CMD if [[ "$(wget -qqO- 'https://duckduckgo.com/?q=what+is+my+ip' | grep -ow 'Your IP address is [0-9.]*[0-9]' | grep -ow '[0-9][0-9.]*')" == "$INITIAL_IP" ]]; then echo "IP address is the same as the non VPN IP address"; exit 1; fi
COPY --from=qmcgaw/dns-trustanchor /named.root /etc/unbound/root.hints
COPY --from=qmcgaw/dns-trustanchor /root.key /etc/unbound/root.key
RUN echo https://dl-3.alpinelinux.org/alpine/v3.8/main > /etc/apk/repositories && \
    apk add -q --progress --no-cache --update openvpn wget ca-certificates iptables unbound && \
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
    chown unbound /etc/unbound/root.key && \
    adduser -S nonrootuser
COPY unbound.conf /etc/unbound/unbound.conf
COPY --from=qmcgaw/malicious-hostnames /malicious-hostnames.bz2 /etc/unbound/malicious-hostnames.bz2
COPY --from=qmcgaw/malicious-ips /malicious-ips.bz2 /etc/unbound/malicious-ips.bz2
COPY entrypoint.sh /entrypoint.sh
RUN chmod 700 /entrypoint.sh
ENTRYPOINT /entrypoint.sh
