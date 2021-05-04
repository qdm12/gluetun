# Start with alpine
FROM alpine:3.13

ENV USER= \
    PASSWORD= \
    REGION="Netherlands" \
    WEBUI_PORT=8888 \
    DNS_SERVERS=209.222.18.222,209.222.18.218


# Download Folder
VOLUME /downloads

# qBittorrent Config Folder
VOLUME /config

# Port for qBittorrent
EXPOSE 8888

ENV DEBIAN_FRONTEND noninteractive

# Ok lets install everything
RUN apk add --no-cache -t .build-deps boost-thread boost-system boost-dev g++ git make automake autoconf libtool libressl-dev qt5-qttools-dev curl unzip dumb-init && \
	apk add --no-cache ca-certificates libressl qt5-qtbase iptables openvpn ack bind-tools python3 && \
	if [ ! -e /usr/bin/python ]; then ln -sf python3 /usr/bin/python ; fi && \
	mkdir /tmp/libtorrent && \
  curl -sSL https://github.com/arvidn/libtorrent/archive/v1.2.13.tar.gz | tar xzC /tmp/libtorrent && \
	cd /tmp/libtorrent/*lib* && \
  ./autotool.sh && \
  ./configure --disable-debug --enable-encryption && \
  make clean && \
  make install && \
	mkdir /tmp/qbittorrent && \
  curl -sSL https://api.github.com/repos/qbittorrent/qBittorrent/tarball/release-4.3.5 | tar xzC /tmp/qbittorrent && \
	cd /tmp/qbittorrent/*qBittorrent* && \
	./configure --disable-gui && \
	make install && \
  mkdir /tmp/openvpn && \
  cd /tmp/openvpn && \
  curl -sSL https://www.privateinternetaccess.com/openvpn/openvpn.zip -o openvpn-nextgen.zip && \
  mkdir -p /openvpn/target && \
  unzip -q openvpn-nextgen.zip -d /openvpn/nextgen && \
  rm *.zip &&  \
  apk del --purge .build-deps && \
	cd / && \
	rm -rf /tmp/* /var/tmp/* /var/cache/apk/* /var/cache/distfiles/* /usr/include/*

COPY ./entrypoint.sh ./qBittorrent.conf /

RUN chmod 500 /entrypoint.sh

# Start point for docker
ENTRYPOINT /entrypoint.sh

# healthcheck
HEALTHCHECK --interval=60s --timeout=15s --start-period=120s \
             CMD curl -LSs 'https://api.ipify.org'
