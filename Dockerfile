FROM alpine
MAINTAINER Quentin McGaw <quentin.mcgaw@gmail.com>
RUN mkdir /pia
WORKDIR /pia
COPY script.sh /pia
RUN apk add --no-cache --update openvpn && \
    apk add --no-cache --update --virtual build-dependencies curl unzip && \
    curl https://www.privateinternetaccess.com/openvpn/openvpn.zip > openvpn.zip && \
    unzip openvpn.zip && rm openvpn.zip && \
    apk del build-dependencies && \
    chmod +x script.sh
ENTRYPOINT ["/pia/script.sh"]