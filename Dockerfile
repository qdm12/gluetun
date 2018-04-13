FROM alpine:3.7
LABEL maintainer="quentin.mcgaw@gmail.com" \
      description="VPN client to private internet access servers using OpenVPN, Alpine and Cloudflare 1.1.1.1 DNS over TLS" \
      download="?MB" \
      size="12.9MB" \
      ram="?MB" \
      cpu_usage="Very low" \
      github="https://github.com/qdm12/private-internet-access-docker"
RUN apk add -q --progress --no-cache --update openvpn unbound && \
    apk add -q --progress --no-cache --update --virtual build-dependencies ca-certificates wget unzip && \
    wget -q https://www.privateinternetaccess.com/openvpn/openvpn.zip && \
    unzip -q openvpn.zip && \
    apk del -q --progress --purge build-dependencies && \
    rm -rf /var/cache/apk/* /etc/unbound/unbound.conf /openvpn.zip
COPY unbound.conf /etc/unbound/unbound.conf
ENTRYPOINT echo "nameserver 127.0.0.1" > /etc/resolv.conf && \
           echo "options ndots:0" >> /etc/resolv.conf && \
           unbound && \
           openvpn --config "$REGION".ovpn --auth-user-pass auth.conf