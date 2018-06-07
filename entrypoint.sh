#!/bin/sh

# Obtaining your original IP address to use for the healthcheck
printf "\nGetting non VPN public IP address..."
export INITIAL_IP=$(wget -qqO- 'https://duckduckgo.com/?q=what+is+my+ip' | grep -ow 'Your IP address is [0-9.]*[0-9]' | grep -ow '[0-9][0-9.]*')
printf "$INITIAL_IP"

# Setting up cloudflare DNS 1.1.1.1 over TLS
printf "\nChanging DNS to localhost..."
echo "nameserver 127.0.0.1" > /etc/resolv.conf
echo "options ndots:0" >> /etc/resolv.conf
printf "DONE"
printf "\nLaunching Unbound daemon to connect to Cloudflare DNS 1.1.1.1 at its TLS endpoint..."
unbound
printf "DONE"
printf "\nSetting firewall for killswitch purposes..."
printf "\n * Detecting local subnet..."
SUBNET=$(ip route show default | tail -n 1 | awk '// {print $1}')
printf "$SUBNET"
printf "\n * Detecting parameters to be used for region $REGION, protocol $PROTOCOL and encryption $ENCRYPTION..."
CONNECTIONSTRING=$(grep -i "/openvpn-$PROTOCOL-$ENCRYPTION/$REGION.ovpn" -e 'privateinternetaccess.com')
PORT=$(echo $CONNECTIONSTRING | cut -d' ' -f3)
PIADOMAIN=$(echo $CONNECTIONSTRING | cut -d' ' -f2)
printf "\n   * Port: $PORT"
printf "\n   * Domain: $PIADOMAIN"
printf "\n     * Detecting IP addresses corresponding to $PIADOMAIN..."
VPNIPS=$(nslookup $PIADOMAIN localhost | tail -n +5 | grep -o '[0-9]\{1,3\}\.[0-9]\{1,3\}\.[0-9]\{1,3\}\.[0-9]\{1,3\}')
for ip in $VPNIPS
do
	printf "\n        $ip"
done
printf "\n * Deleting all iptables rules..."
iptables --flush
iptables --delete-chain
iptables -t nat --flush
iptables -t nat --delete-chain
iptables -P OUTPUT DROP
printf "DONE"
printf "\n * Adding rules to accept local loopback traffic..."
iptables -A INPUT -j ACCEPT -i lo
iptables -A OUTPUT -j ACCEPT -o lo
printf "DONE"
printf "\n * Adding rules to accept traffic of subnet $SUBNET..."
#iptables -A INPUT --src $SUBNET -j ACCEPT -i eth0
iptables -A OUTPUT -d $SUBNET -j ACCEPT -o eth0
printf "DONE"
for ip in $VPNIPS
do
	printf "\n * Adding rules to accept traffic with $ip on port $PROTOCOL $PORT..."
	iptables -A OUTPUT -j ACCEPT -d $ip -o eth0 -p $PROTOCOL -m $PROTOCOL --dport $PORT
	iptables -A INPUT -j ACCEPT -s $ip -i eth0 -p $PROTOCOL -m $PROTOCOL --sport $PORT
	printf "DONE"
done
printf "\n * Adding rules to accept traffic going through the tun device..."
iptables -A INPUT -j ACCEPT -i tun0
iptables -A OUTPUT -j ACCEPT -o tun0
printf "DONE"
printf "\n * Starting OpenVPN using the following parameters:"
printf "\n   * Domain: $PIADOMAIN"
printf "\n   * Port: $PORT"
printf "\n   * Protocol: $PROTOCOL"
printf "\n   * Encryption: $ENCRYPTION"
cd /openvpn-$PROTOCOL-$ENCRYPTION
openvpn --config "$REGION.ovpn" --auth-user-pass /auth.conf
printf "\nExiting...\n\n"