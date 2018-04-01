FROM alpine:3.7
LABEL maintainer="quentin.mcgaw@gmail.com" \
      description="VPN client container to private internet access servers based on Alpine Linux and OpenVPN" \
      download="3.3MB" \
      size="8MB" \
      ram="4.3MB" \
      cpu_usage="Very low"
      github="https://github.com/qdm12/private-internet-access-docker"
COPY script.sh .
RUN chmod +x script.sh && \
    apk add -q --progress --no-cache --update openvpn && \
    apk add -q --progress --no-cache --update --virtual build-dependencies wget unzip && \
    wget https://www.privateinternetaccess.com/openvpn/openvpn.zip && \
    unzip openvpn.zip && \
    rm openvpn.zip && \
    apk del -q --progress --purge build-dependencies && \
    rm -rf /var/cache/apk/*
ENTRYPOINT /script.sh