#!/bin/sh

printf "\nDetecting details from public IP address..."
export CITY=$(wget -qO- -T 2 https://ipinfo.io/city)
export ORG=$(wget -qO- -T 2 https://ipinfo.io/org)
printf "DONE\nOrganization: $ORG\nCountry: $COUNTRY\nCity: $CITY\nChanging DNS to localhost..."
echo "nameserver 127.0.0.1" > /etc/resolv.conf
echo "options ndots:0" >> /etc/resolv.conf
printf "DONE\nStarting Unbound to connect to Cloudflare DNS 1.1.1.1 at its TLS endpoint..."
unbound
printf "DONE\nStarting OpenVPN using $PROTOCOL with $ENCRYPTION encryption\n"
DIR=/openvpn-$PROTOCOL-$ENCRYPTION
openvpn --config $DIR/$REGION.ovpn --auth-user-pass /auth.conf --ca $DIR/ca.rsa.*.crt --crl-verify $DIR/ca.rsa.*.crt
printf "\n\nExiting..."
