# Start with alpine
FROM alpine:3.15

ENV USER= \
    PASSWORD= \
    REGION="Netherlands" \
    WEBUI_PORT=8888 \
    DNS_SERVERS=209.222.18.222,209.222.18.218,103.196.38.38,103.196.38.39


# Download Folder
VOLUME /downloads

# qBittorrent Config Folder
VOLUME /config

# Port for qBittorrent
EXPOSE 8888

ENV DEBIAN_FRONTEND noninteractive

# Ok lets install everything
RUN apk add --no-cache -t .build-deps autoconf automake build-base cmake curl git libtool linux-headers perl pkgconf python3-dev re2c tar unzip icu-dev libexecinfo-dev openssl-dev qt5-qtbase-dev qt5-qttools-dev zlib-dev qt5-qtsvg-dev && \
	apk add --no-cache ca-certificates libexecinfo libressl qt5-qtbase iptables openvpn ack bind-tools python3 && \
	if [ ! -e /usr/bin/python ]; then ln -sf python3 /usr/bin/python ; fi && \
  curl -sNLk --retry 5 https://boostorg.jfrog.io/artifactory/main/release/1.76.0/source/boost_1_76_0.tar.gz | tar xzC /tmp && \
  curl -sSL --retry 5 https://github.com/ninja-build/ninja/archive/refs/tags/v1.10.2.tar.gz | tar xzC /tmp && \
	cd /tmp/*ninja* && \
  cmake -Wno-dev -B build \
  	-D CMAKE_CXX_STANDARD=17 \
  	-D CMAKE_INSTALL_PREFIX="/usr/local" && \
  cmake --build build && \
  cmake --install build && \
  curl -sSL --retry 5 https://github.com/arvidn/libtorrent/archive/v1.2.15.tar.gz | tar xzC /tmp && \
	cd /tmp/*libtorrent* && \
  cmake -Wno-dev -G Ninja -B build \
    -D CMAKE_BUILD_TYPE="Release" \
    -D CMAKE_CXX_STANDARD=17 \
    -D BOOST_INCLUDEDIR="/tmp/boost_1_76_0/" \
    -D CMAKE_INSTALL_LIBDIR="lib" \
    -D CMAKE_INSTALL_PREFIX="/usr/local" && \
  cmake --build build && \
  cmake --install build && \
  curl -sSL --retry 5 https://api.github.com/repos/qbittorrent/qBittorrent/tarball/release-4.4.1 | tar xzC /tmp && \
	cd /tmp/*qBittorrent* && \
  cmake -Wno-dev -G Ninja -B build \
    -D CMAKE_BUILD_TYPE="release" \
    -D GUI=OFF \
    -D CMAKE_CXX_STANDARD=17 \
    -D BOOST_INCLUDEDIR="/tmp/boost_1_76_0/" \
    -D CMAKE_CXX_STANDARD_LIBRARIES="/usr/lib/libexecinfo.so" \
    -D CMAKE_INSTALL_PREFIX="/usr/local" && \
  cmake --build build && \
  cmake --install build && \
  mkdir /tmp/openvpn && \
  cd /tmp/openvpn && \
  curl -sSL --retry 5 https://www.privateinternetaccess.com/openvpn/openvpn.zip -o openvpn-nextgen.zip && \
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
