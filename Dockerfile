FROM alpine:3.8 AS rootanchors
RUN apk add -q --progress wget perl-xml-xpath && \
    wget -q https://www.internic.net/domain/named.root -O named.root && \
    echo "602f28581292bf5e50c8137c955173e6  named.root" > hashes.md5 && \
    md5sum -c hashes.md5 && \
    wget -q https://data.iana.org/root-anchors/root-anchors.xml -O root-anchors.xml && \
    echo "1b2a628d1ff22d4dc7645cfc89f21b6a575526439c6706ecf853e6fff7099dc8  root-anchors.xml" > hashes.sha256 && \
    sha256sum -c hashes.sha256 && \
    KEYTAGS=$(xpath -q -e '/TrustAnchor/KeyDigest/KeyTag/node()' root-anchors.xml) && \
    ALGORITHMS=$(xpath -q -e '/TrustAnchor/KeyDigest/Algorithm/node()' root-anchors.xml) && \
    DIGESTTYPES=$(xpath -q -e '/TrustAnchor/KeyDigest/DigestType/node()' root-anchors.xml) && \
    DIGESTS=$(xpath -q -e '/TrustAnchor/KeyDigest/Digest/node()' root-anchors.xml) && \
    i=1 && \
    while [ 1 ]; do \
      KEYTAG=$(echo $KEYTAGS | cut -d" " -f$i); \
      [ "$KEYTAG" != "" ] || break; \
      ALGORITHM=$(echo $ALGORITHMS | cut -d" " -f$i); \
      DIGESTTYPE=$(echo $DIGESTTYPES | cut -d" " -f$i); \
      DIGEST=$(echo $DIGESTS | cut -d" " -f$i); \
      echo ". IN DS $KEYTAG $ALGORITHM $DIGESTTYPE $DIGEST" >> /root.key; \
      i=`expr $i + 1`; \
    done;

FROM alpine:3.8 AS blocks
RUN apk add -q --progress wget && \
    wget -q https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts -O temp && \
    sed -i '/\(^[ \|\t]*#\)\|\(^[ ]\+\)\|\(^$\)\|\(^[\n\|\r\|\r\n][ \|\t]*$\)\|\(^127.0.0.1\)\|\(^255.255.255.255\)\|\(^::1\)\|\(^fe80\)\|\(^ff00\)\|\(^ff02\)\|\(^0.0.0.0 0.0.0.0\)/d' temp && \
    sed -i 's/\([ \|\t]*#.*$\)\|\(\r\)\|\(0.0.0.0 \)//g' temp && \
    cat temp >> allHostnames && \
    wget -q https://raw.githubusercontent.com/CHEF-KOCH/NSABlocklist/master/HOSTS -O temp && \
    sed -i '/\(^[ \|\t]*#\)\|\(^[ ]\+\)\|\(^$\)\|\(^[\n\|\r\|\r\n][ \|\t]*$\)\|\(^127.0.0.1\)/d' temp && \
    sed -i 's/\([ \|\t]*#.*$\)\|\(\r\)\|\(0.0.0.0 \)//g' temp && \
    cat temp >> allHostnames && \
    wget -q https://raw.githubusercontent.com/k0nsl/unbound-blocklist/master/blocks.conf -O temp && \
    sed -i '/\(^[ \|\t]*#\)\|\(^[ ]\+\)\|\(^$\)\|\(^[\n\|\r\|\r\n][ \|\t]*$\)\|\(^local-data\)/d' temp && \
    sed -i 's/\([ \|\t]*#.*$\)\|\(\r\)\|\(local-zone: \"\)\|\(\" redirect\)//g' temp && \
    cat temp >> allHostnames && \
    wget -q https://raw.githubusercontent.com/notracking/hosts-blocklists/master/domains.txt -O temp && \
    sed -i '/\(^[ \|\t]*#\)\|\(^[ ]\+\)\|\(^$\)\|\(^[\n\|\r\|\r\n][ \|\t]*$\)\|\(::$\)/d' temp && \
    sed -i 's/\([ \|\t]*#.*$\)\|\(\r\)\|\(address=\/\)\|\(\/0.0.0.0$\)//g' temp && \
    cat temp >> allHostnames && \
    wget -q https://raw.githubusercontent.com/notracking/hosts-blocklists/master/hostnames.txt -O temp && \
    sed -i '/\(^[ \|\t]*#\)\|\(^[ ]\+\)\|\(^$\)\|\(^[\n\|\r\|\r\n][ \|\t]*$\)\|\(^::\)/d' temp && \
    sed -i 's/\([ \|\t]*#.*$\)\|\(\r\)\|\(^0.0.0.0 \)//g' temp && \
    cat temp >> allHostnames && \
    sort -o allHostnames allHostnames && \
    cat allHostnames | uniq -i -u > uniqueHostnames && \
    cat allHostnames | uniq -i -d > duplicateHostnames && \
    duplicates=$(($(wc -l < allHostnames)-$(wc -l < uniqueHostnames)-$(wc -l < duplicateHostnames))) && \
    echo "Removed $duplicates duplicates in a total of $(wc -l < allHostnames) hostnames" && \
    mv uniqueHostnames allHostnames && \
    cat duplicateHostnames >> allHostnames && \
    sort -o allHostnames allHostnames && \
    sed -i '/\(psma01.com.\)\|\(psma02.com.\)\|\(psma03.com.\)\|\(MEZIAMUSSUCEMAQUEUE.SU\)/d' allHostnames && \
    while read line; do printf "local-zone: \"$line\" static\n" >> blocks-malicious.conf; done < allHostnames && \
    tar -cjf blocks-malicious.conf.bz2 blocks-malicious.conf

FROM alpine:3.8
LABEL maintainer="quentin.mcgaw@gmail.com" \
      description="VPN client to private internet access servers using OpenVPN, IPtables firewall, DNS over TLS with Unbound and Alpine Linux" \
      download="6.6MB" \
      size="15.7MB" \
      ram="13MB" \
      cpu_usage="Low" \
      github="https://github.com/qdm12/private-internet-access-docker"
COPY --from=rootanchors /named.root /etc/unbound/root.hints
COPY --from=rootanchors /root.key /etc/unbound/root.key
COPY --from=blocks /blocks-malicious.conf.bz2 /etc/unbound/blocks-malicious.conf.bz2
HEALTHCHECK --interval=5m --timeout=15s --start-period=10s --retries=2 \
            CMD if [[ "$(wget -qqO- 'https://duckduckgo.com/?q=what+is+my+ip' | grep -ow 'Your IP address is [0-9.]*[0-9]' | grep -ow '[0-9][0-9.]*')" == "$INITIAL_IP" ]]; then echo "IP address is the same as the non VPN IP address"; exit 1; fi
ENV ENCRYPTION=strong \
    PROTOCOL=tcp \
    REGION=Germany \
    BLOCK_MALICIOUS=off
RUN apk add -q --progress --no-cache --update openvpn ca-certificates iptables unbound && \
    apk add -q --progress --no-cache --update --virtual=build-dependencies unzip && \
    mkdir /openvpn-udp-normal /openvpn-udp-strong /openvpn-tcp-normal /openvpn-tcp-strong && \
    wget -q https://www.privateinternetaccess.com/openvpn/openvpn.zip \
            https://www.privateinternetaccess.com/openvpn/openvpn-strong.zip \
            https://www.privateinternetaccess.com/openvpn/openvpn-tcp.zip \
            https://www.privateinternetaccess.com/openvpn/openvpn-strong-tcp.zip && \
    unzip -q openvpn.zip -d /openvpn-udp-normal && \
    unzip -q openvpn-strong.zip -d /openvpn-udp-strong && \
    unzip -q openvpn-tcp.zip -d /openvpn-tcp-normal && \
    unzip -q openvpn-strong-tcp.zip -d /openvpn-tcp-strong && \
    apk del -q --progress --purge build-dependencies && \
    rm -rf /*.zip /var/cache/apk/* /etc/unbound/unbound.conf && \
    chown unbound /etc/unbound/root.key && \
    adduser -S nonrootuser
COPY unbound.conf /etc/unbound/unbound.conf
COPY entrypoint.sh /
RUN chmod +x /entrypoint.sh
ENTRYPOINT /entrypoint.sh
