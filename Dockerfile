FROM alpine
MAINTAINER Quentin McGaw <quentin.mcgaw@gmail.com>
RUN mkdir /pia && cd /pia
COPY script.sh /pia
RUN apk add --no-cache openvpn && \
    apk add --no-cache --virtual build-dependencies curl unzip && \
    curl https://www.privateinternetaccess.com/openvpn/openvpn.zip > openvpn.zip && \
    unzip openvpn.zip && rm openvpn.zip && \
    apk del build-dependencies && \
    chmod +x script.sh
ENTRYPOINT ["/pia/script.sh"]