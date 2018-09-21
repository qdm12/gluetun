#!/bin/sh

printf "=== PIA CONTAINER ==="

############################################
# SETTING DNS OVER TLS TO 1.1.1.1 / 1.0.0.1
############################################
printf "\nChanging DNS to localhost..."
echo "nameserver 127.0.0.1" > /etc/resolv.conf
echo "options ndots:0" >> /etc/resolv.conf
printf "DONE"
printf "\nLaunching Unbound daemon to connect to Cloudflare DNS 1.1.1.1 at its TLS endpoint..."
unbound
printf "DONE"

############################################
# ORIGINAL IP FOR HEALTHCHECK
############################################
printf "\nGetting non VPN public IP address..."
export INITIAL_IP=$(wget -qqO- 'https://duckduckgo.com/?q=what+is+my+ip' | grep -ow 'Your IP address is [0-9.]*[0-9]' | grep -ow '[0-9][0-9.]*')
printf "$INITIAL_IP"

############################################
# FIREWALL
############################################
printf "\nSetting firewall for killswitch purposes..."
printf "\n * Detecting local subnet..."
SUBNET=$(ip route show default | tail -n 1 | awk '// {print $1}')
printf "$SUBNET"
printf "\n * Reading parameters to be used for region $REGION, protocol $PROTOCOL and encryption $ENCRYPTION..."
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
ip6tables --flush
ip6tables --delete-chain
printf "DONE"
iptables -F OUTPUT
iptables -P OUTPUT DROP
ip6tables -F OUTPUT 2>/dev/null
ip6tables -P OUTPUT DROP 2>/dev/null
printf "\n * Adding rules to accept local loopback traffic..."
iptables -A OUTPUT -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT
iptables -A OUTPUT -o lo -j ACCEPT
iptables -A INPUT -i lo -j ACCEPT
ip6tables -A OUTPUT -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT 2>/dev/null
ip6tables -A OUTPUT -o lo -j ACCEPT 2>/dev/null
ip6tables -A INPUT -i lo -j ACCEPT 2>/dev/null
printf "DONE"
printf "\n * Adding rules to accept traffic of subnet $SUBNET..."
iptables -A OUTPUT -d $SUBNET -j ACCEPT
ip6tables -A OUTPUT -d $SUBNET -j ACCEPT 2>/dev/null
printf "DONE"
for ip in $VPNIPS
do
    printf "\n * Adding rules to accept traffic with $ip on port $PROTOCOL $PORT..."
    iptables -A OUTPUT -j ACCEPT -d $ip -o eth0 -p $PROTOCOL -m $PROTOCOL --dport $PORT
    iptables -A INPUT -j ACCEPT -s $ip -i eth0 -p $PROTOCOL -m $PROTOCOL --sport $PORT
    ip6tables -A OUTPUT -j ACCEPT -d $ip -o eth0 -p $PROTOCOL -m $PROTOCOL --dport $PORT 2>/dev/null
    ip6tables -A INPUT -j ACCEPT -s $ip -i eth0 -p $PROTOCOL -m $PROTOCOL --sport $PORT 2>/dev/null
    printf "DONE"
done
printf "\n * Adding rules to accept traffic going through the tun device..."
iptables -A OUTPUT -o tun0 -j ACCEPT
iptables -A INPUT -i tun0 -j ACCEPT
ip6tables -A OUTPUT -o tun0 -j ACCEPT 2>/dev/null
ip6tables -A INPUT -i tun0 -j ACCEPT 2>/dev/null
printf "DONE"
printf "\n * Allowing outgoing DNS queries on port 53 UDP..."
iptables -A OUTPUT -p udp -m udp --dport 53 -j ACCEPT
ip6tables -A OUTPUT -p udp -m udp --dport 53 -j ACCEPT 2>/dev/null
printf "DONE"

############################################
# SUMMARY
############################################
printf "\nStarting OpenVPN using the following parameters:"
printf "\n * Domain: $PIADOMAIN"
printf "\n * Port: $PORT"
printf "\n * Protocol: $PROTOCOL"
printf "\n * Encryption: $ENCRYPTION"

############################################
# OPENVPN LAUNCH
############################################
cd /openvpn-$PROTOCOL-$ENCRYPTION
printf "\nSwitching from root to nonrootuser..."
su -l nonrootuser
printf "DONE\n"
openvpn --config "$REGION.ovpn" --auth-user-pass /auth.conf

############################################
# CLEANUP
############################################
printf "\nExiting...\n\n"