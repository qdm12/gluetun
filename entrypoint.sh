#!/bin/sh

printf "\nGetting public IP address..."
export INITIAL_IP=$(wget -qqO- 'https://duckduckgo.com/?q=what+is+my+ip' | grep -ow 'Your IP address is [0-9.]*[0-9]' | grep -ow '[0-9][0-9.]*')
printf "DONE\nChanging DNS to localhost..."
echo "nameserver 127.0.0.1" > /etc/resolv.conf
echo "options ndots:0" >> /etc/resolv.conf
printf "DONE\nStarting Unbound to connect to Cloudflare DNS 1.1.1.1 at its TLS endpoint..."
unbound
printf "DONE\nSetting firewall for killswitch purposes...\n  Detecting local subnet..."
SUBNET=$(ip route show default | tail -n 1 | awk '// {print $1}')
printf "$SUBNET\n  Detecting IP addresses corresponding to $REGION.privateinternetaccess.com..."
VPNIPS=$(nslookup $REGION.privateinternetaccess.com localhost | tail -n +5 | grep -o '[0-9]\{1,3\}\.[0-9]\{1,3\}\.[0-9]\{1,3\}\.[0-9]\{1,3\}')
for ip in $VPNIPS
do
	printf "\n    $ip"
done
printf "\n  Deleting all iptables rules..."
iptables --flush
iptables --delete-chain
iptables -t nat --flush
iptables -t nat --delete-chain
iptables -P OUTPUT DROP
printf "DONE\n  Adding rules to accept local loopback traffic..."
iptables -A INPUT -j ACCEPT -i lo
iptables -A OUTPUT -j ACCEPT -o lo
printf "DONE\n  Adding rules to accept traffic of subnet $SUBNET..."
#iptables -A INPUT --src $SUBNET -j ACCEPT -i eth0
iptables -A OUTPUT -d $SUBNET -j ACCEPT -o eth0
printf "DONE\n  Determining port to be used with PIA..."
if [ "$PROTOCOL-$ENCRYPTION" == "tcp-normal" ]; then
	PORT=502
elif [ "$PROTOCOL-$ENCRYPTION" == "tcp-strong" ]; then
	PORT=501
elif [ "$PROTOCOL-$ENCRYPTION" == "udp-normal" ]; then
	PORT=1198
elif [ "$PROTOCOL-$ENCRYPTION" == "udp-strong" ]; then
	PORT=1197
fi
printf "$PROTOCOL $PORT"
for ip in $VPNIPS
do
	printf "\n  Adding rules to accept traffic with VPN IP address $ip on port $PROTOCOL $PORT..."
	iptables -A OUTPUT -j ACCEPT -d $ip -o eth0 -p $PROTOCOL -m $PROTOCOL --dport $PORT
	iptables -A INPUT -j ACCEPT -s $ip -i eth0 -p $PROTOCOL -m $PROTOCOL --sport $PORT
	printf "DONE"
done
printf "\n  Adding rules to accept traffic going through the tun device..."
iptables -A INPUT -j ACCEPT -i tun0
iptables -A OUTPUT -j ACCEPT -o tun0
printf "DONE\nStarting OpenVPN using $PROTOCOL with $ENCRYPTION encryption\n"
cd /openvpn-$PROTOCOL-$ENCRYPTION
openvpn --config "$REGION.ovpn" --auth-user-pass /auth.conf
printf "\n\nExiting..."