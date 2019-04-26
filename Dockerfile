ARG BASE_IMAGE=alpine
ARG ALPINE_VERSION=3.9

FROM ${BASE_IMAGE}:${ALPINE_VERSION}
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
      org.label-schema.docker.cmd="docker run -d --cap-add=NET_ADMIN --device=/dev/net/tun -e USER=js89ds7 -e PASSWORD=8fd9s239G qmcgaw/private-internet-access" \
      org.label-schema.docker.cmd.devel="docker run -it --rm --cap-add=NET_ADMIN --device=/dev/net/tun -e USER=js89ds7 -e PASSWORD=8fd9s239G qmcgaw/private-internet-access" \
      org.label-schema.docker.params="REGION=PIA region,PROTOCOL=udp/tcp,ENCRYPTION=strong/normal,BLOCK_MALICIOUS=on/off,BLOCK_NSA=on/off,UNBLOCK=allowed hostnames,USER=PIA user,PASSWORD=PIA password,EXTRA_SUBNETS=extra subnets to allow on the firewall,NONROOT=yes/no" \
      org.label-schema.version="" \
      image-size="19.6MB" \
      ram-usage="13MB to 80MB" \
      cpu-usage="Low to Medium"
ENV USER= \
    PASSWORD= \
    ENCRYPTION=strong \
    PROTOCOL=udp \
    REGION="CA Montreal" \
    BLOCK_MALICIOUS=off \
    BLOCK_NSA=off \
    UNBLOCK= \
    EXTRA_SUBNETS= \
    NONROOT=no
ENTRYPOINT /entrypoint.sh
HEALTHCHECK --interval=3m --timeout=3s --start-period=20s --retries=1 CMD /healthcheck.sh
RUN apk add -q --progress --no-cache --update openvpn wget ca-certificates iptables unbound unzip && \
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
    rm -rf /*.zip /var/cache/apk/* /etc/unbound/unbound.conf /usr/sbin/unbound-anchor /usr/sbin/unbound-checkconf /usr/sbin/unbound-control /usr/sbin/unbound-control-setup /usr/sbin/unbound-host && \
    adduser nonrootuser -D -H --uid 1000 && \
    wget -q https://raw.githubusercontent.com/qdm12/updated/master/files/named.root.updated -O /etc/unbound/root.hints && \
    wget -q https://raw.githubusercontent.com/qdm12/updated/master/files/root.key.updated -O /etc/unbound/root.key && \
    cd /tmp && \
    wget -q https://raw.githubusercontent.com/qdm12/updated/master/files/malicious-hostnames.updated -O malicious-hostnames && \
    wget -q https://raw.githubusercontent.com/qdm12/updated/master/files/nsa-hostnames.updated -O nsa-hostnames && \
    wget -q https://raw.githubusercontent.com/qdm12/updated/master/files/malicious-ips.updated -O malicious-ips && \
    while read hostname; do echo "local-zone: \""$hostname"\" static" >> blocks-malicious.conf; done < malicious-hostnames && \
    while read ip; do echo "private-address: $ip" >> blocks-malicious.conf; done < malicious-ips && \
    tar -cjf /etc/unbound/blocks-malicious.bz2 blocks-malicious.conf && \
    while read hostname; do echo "local-zone: \""$hostname"\" static" >> blocks-nsa.conf; done < nsa-hostnames && \
    tar -cjf /etc/unbound/blocks-nsa.bz2 blocks-nsa.conf && \
    rm -f /tmp/*
COPY unbound.conf /etc/unbound/unbound.conf
COPY entrypoint.sh healthcheck.sh /
RUN chown nonrootuser -R /etc/unbound && \
    chmod 700 /etc/unbound && \
    chmod 500 /entrypoint.sh healthcheck.sh && \
    chmod 400 \
        /etc/unbound/root.hints \
        /etc/unbound/root.key \
        /etc/unbound/unbound.conf \
        /etc/unbound/*.bz2
