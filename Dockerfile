# Latest version of ubuntu
FROM ubuntu:18.04

ENV USER= \
    PASSWORD= \
    ENCRYPTION=strong \
    PROTOCOL=udp \
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
RUN apt-get update && \
    apt-get install --no-install-recommends -y apt-utils software-properties-common && \
    apt-get install --no-install-recommends -y unzip build-essential pkg-config automake libtool git zlib1g-dev libboost-dev libboost-system-dev libboost-chrono-dev libboost-random-dev libssl-dev libgeoip-dev curl cmake qtbase5-dev qttools5-dev-tools libqt5svg5-dev && \
	apt-get install --no-install-recommends -y ca-certificates openvpn openvpn-systemd-resolved wget ca-certificates iptables dnsutils iputils-ping net-tools ack && \
	LIBTORRENT_URL=$(curl -sSL https://api.github.com/repos/arvidn/libtorrent/tags | grep tarball_url | head -n 1 | cut -d '"' -f 4) && \
	mkdir /tmp/libtorrent && \
	curl -sSL https://api.github.com/repos/arvidn/libtorrent/tarball/libtorrent_1_1_12 | tar xzC /tmp/libtorrent && \
	cd /tmp/libtorrent/*lib* && \
	mkdir build && \
	cd build && \
	cmake .. && \
	make install && \
	mkdir /tmp/qbittorrent && \
	curl -sSL https://api.github.com/repos/qbittorrent/qBittorrent/tarball/release-4.2.1 | tar xzC /tmp/qbittorrent && \
	cd /tmp/qbittorrent/*qbittorrent* && \
	./configure --disable-gui CXXFLAGS="-std=c++14" && \
	make -j$(nproc) && \
	make install && \
    wget -q https://www.privateinternetaccess.com/openvpn/openvpn.zip \
    https://www.privateinternetaccess.com/openvpn/openvpn-strong.zip \
    https://www.privateinternetaccess.com/openvpn/openvpn-tcp.zip \
    https://www.privateinternetaccess.com/openvpn/openvpn-strong-tcp.zip && \
    mkdir -p /openvpn/target && \
    unzip -q openvpn.zip -d /openvpn/udp-normal && \
    unzip -q openvpn-strong.zip -d /openvpn/udp-strong && \
    unzip -q openvpn-tcp.zip -d /openvpn/tcp-normal && \
    unzip -q openvpn-strong-tcp.zip -d /openvpn/tcp-strong && \
    apt-get purge -y -qq unzip software-properties-common wget apt-utils build-essential pkg-config automake libtool git zlib1g-dev libboost-dev libboost-system-dev libboost-chrono-dev libboost-random-dev libssl-dev libgeoip-dev curl cmake qtbase5-dev qttools5-dev-tools libqt5svg5-dev && \
    apt-get clean -qq && \
    apt-get autoclean -qq && \
    rm -rf /*.zip /tmp/* /var/tmp/* /var/lib/apt/lists/*


COPY entrypoint.sh qBittorrent.conf /

RUN chmod 500 /entrypoint.sh
