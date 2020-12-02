# Start with alpine
FROM alpine:3.12.0

ENV USER= \
    PASSWORD= \
    REGION="Netherlands" \
    WEBUI_PORT=8888 \
    DNS_SERVERS=209.222.18.222,209.222.18.218


# Start point for docker
ENTRYPOINT /entrypoint.sh

# Download Folder
VOLUME /downloads

# qBittorrent Config Folder
VOLUME /config

# Port for qBittorrent
EXPOSE 8888

ENV DEBIAN_FRONTEND noninteractive

# Ok lets install everything
RUN apk add --no-cache -t .build-deps boost-thread boost-system boost-dev g++ git make cmake libressl-dev qt5-qttools-dev curl unzip dumb-init && \
	apk add --no-cache ca-certificates libressl qt5-qtbase iptables openvpn ack bind-tools python3 && \
	if [ ! -e /usr/bin/python ]; then ln -sf python3 /usr/bin/python ; fi && \
	LIBTORRENT_URL=$(curl -sSL https://api.github.com/repos/arvidn/libtorrent/tags | grep tarball_url | head -n 1 | cut -d '"' -f 4) && \
	mkdir /tmp/libtorrent && \
  curl -sSL https://api.github.com/repos/arvidn/libtorrent/tarball/libtorrent_1_2_7 | tar xzC /tmp/libtorrent && \
	cd /tmp/libtorrent/*lib* && \
	mkdir -p cmake-build-dir/release && \
	cd cmake-build-dir/release && \
	cmake -DCMAKE_BUILD_TYPE=Release -DCMAKE_CXX_STANDARD=14 -G "Unix Makefiles" ../.. && \
	make install && \
	mkdir /tmp/qbittorrent && \
	curl -sSL https://api.github.com/repos/qbittorrent/qBittorrent/tarball/release-4.3.1 | tar xzC /tmp/qbittorrent && \
	cd /tmp/qbittorrent/*qbittorrent* && \
	PKG_CONFIG_PATH=/usr/local/lib64/pkgconfig ./configure --disable-gui && \
	make install && \
	export LD_LIBRARY_PATH=/usr/local/lib:/usr/local/lib64:${LD_LIBRARY_PATH} && \
	mkdir /tmp/openvpn && \
	cd /tmp/openvpn && \
	curl -sSL https://www.privateinternetaccess.com/openvpn/openvpn.zip -o openvpn-nextgen.zip && \
	mkdir -p /openvpn/target && \
	unzip -q openvpn-nextgen.zip -d /openvpn/nextgen && \
	apk del --purge .build-deps && \
	cd / && \
	rm -rf /tmp/* /var/tmp/* /var/cache/apk/* /var/cache/distfiles/* /usr/include/*



COPY entrypoint.sh qBittorrent.conf /

RUN chmod 500 /entrypoint.sh
