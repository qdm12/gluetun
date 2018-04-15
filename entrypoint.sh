#!/bin/sh

printf "Changing DNS to localhost..."
echo "nameserver 127.0.0.1" > /etc/resolv.conf
echo "options ndots:0" >> /etc/resolv.conf
printf "DONE\nStarting Unbound to connect to Cloudflare DNS 1.1.1.1 at its TLS endpoint TCP 853..."
unbound
printf "DONE\nStarting OpenVPN using $PROTOCOL with $ENCRYPTION encryption\n"
DIR=/openvpn-$PROTOCOL-$ENCRYPTION
openvpn --config $DIR/$REGION.ovpn --auth-user-pass /auth.conf --ca $DIR/ca.rsa.*.crt --crl-verify $DIR/ca.rsa.*.crt
printf "\n\nExiting..."
