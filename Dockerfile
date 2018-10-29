ARG ALPINE_VERSION=3.8

FROM alpine:${ALPINE_VERSION}
ARG BUILD_DATE
ARG VCS_REF
LABEL org.label-schema.schema-version="1.0.0-rc1" \
      maintainer="quentin.mcgaw@gmail.com" \
      org.label-schema.build-date=$BUILD_DATE \
      org.label-schema.vcs-ref=$VCS_REF \
      org.label-schema.vcs-url="https://github.com/qdm12/private-internet-access-docker" \
      org.label-schema.url="https://github.com/qdm12/private-internet-access-docker" \
      org.label-schema.vcs-description="VPN client to tunnel to private internet access servers using OpenVPN, IPtables, DNS over TLS and Alpine Linux" \
      org.label-schema.vcs-usage="https://github.com/qdm12/private-internet-access-docker/blob/master/README.md#setup" \
      org.label-schema.docker.cmd="docker run -d -v ./auth.conf:/auth.conf:ro --cap-add=NET_ADMIN --device=/dev/net/tun qmcgaw/private-internet-access" \
      org.label-schema.docker.cmd.devel="docker run -it --rm -v ./auth.conf:/auth.conf:ro --cap-add=NET_ADMIN --device=/dev/net/tun qmcgaw/private-internet-access" \
      org.label-schema.docker.params="" \
      org.label-schema.version="" \
      image-size="17.1MB" \
      ram-usage="13MB to 80MB" \
      cpu-usage="Low"
ENV ENCRYPTION=strong \
    PROTOCOL=tcp \
    REGION="CA Montreal" \
    BLOCK_MALICIOUS=off
HEALTHCHECK --interval=5m --timeout=15s --start-period=10s --retries=2 \
            CMD if [[ "$(wget -qqO- 'https://duckduckgo.com/?q=what+is+my+ip' | grep -ow 'Your IP address is [0-9.]*[0-9]' | grep -ow '[0-9][0-9.]*')" == "$INITIAL_IP" ]]; then echo "IP address is the same as the non VPN IP address"; exit 1; fi
RUN V_ALPINE="v$(cat /etc/alpine-release | grep -oE '[0-9]+\.[0-9]+')" && \
    echo https://dl-3.alpinelinux.org/alpine/$V_ALPINE/main > /etc/apk/repositories && \
    apk add -q --progress --no-cache --update openvpn wget ca-certificates iptables unbound unzip && \
    wget -q https://www.privateinternetaccess.com/openvpn/openvpn.zip \
            https://www.privateinternetaccess.com/openvpn/openvpn-strong.zip \
            https://www.privateinternetaccess.com/openvpn/openvpn-tcp.zip \
            https://www.privateinternetaccess.com/openvpn/openvpn-strong-tcp.zip && \
    unzip -q openvpn.zip -d /openvpn-udp-normal && \
    unzip -q openvpn-strong.zip -d /openvpn-udp-strong && \
    unzip -q openvpn-tcp.zip -d /openvpn-tcp-normal && \
    unzip -q openvpn-strong-tcp.zip -d /openvpn-tcp-strong && \
    apk del -q --progress --purge unzip && \
    rm -rf /*.zip /var/cache/apk/* /etc/unbound/unbound.conf && \
    adduser -S nonrootuser
COPY --from=qmcgaw/dns-trustanchor /named.root /etc/unbound/root.hints
COPY --from=qmcgaw/dns-trustanchor /root.key /etc/unbound/root.key
COPY --from=qmcgaw/malicious-hostnames /malicious-hostnames.bz2 /etc/unbound/malicious-hostnames.bz2
COPY --from=qmcgaw/malicious-ips /malicious-ips.bz2 /etc/unbound/malicious-ips.bz2
COPY unbound.conf /etc/unbound/unbound.conf
COPY entrypoint.sh /entrypoint.sh
RUN chown unbound /etc/unbound/root.key && \
    chmod 700 /entrypoint.sh
ENTRYPOINT /entrypoint.sh
