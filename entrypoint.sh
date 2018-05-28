#!/bin/sh

printf "\nGetting public IP address..."
export INITIAL_IP=$(wget -qqO- 'https://duckduckgo.com/?q=what+is+my+ip' | grep -ow 'Your IP address is [0-9.]*[0-9]' | grep -ow '[0-9][0-9.]*')
printf "DONE\nChanging DNS to localhost..."
echo "nameserver 127.0.0.1" > /etc/resolv.conf
echo "options ndots:0" >> /etc/resolv.conf
printf "DONE\nStarting Unbound to connect to Cloudflare DNS 1.1.1.1 at its TLS endpoint..."
unbound
printf "DONE\nStarting OpenVPN using $PROTOCOL with $ENCRYPTION encryption\n"
cd /openvpn-$PROTOCOL-$ENCRYPTION
openvpn --config "$REGION.ovpn" --auth-user-pass /auth.conf
printf "\n\nExiting..."
