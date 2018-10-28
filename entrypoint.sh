#!/bin/sh

printf "\n ========================================="
printf "\n ========================================="
printf "\n ============= PIA CONTAINER ============="
printf "\n ========================================="
printf "\n ========================================="
printf "\n == by github.com/qdm12 - Quentin McGaw ==\n"

printf "\nOpenVPN version: $(openvpn --version | head -n 1 | grep -oE "OpenVPN [0-9\.]* " | cut -d" " -f2)"
printf "\nUnbound version: $(unbound -h | grep "Version" | cut -d" " -f2)"
printf "\nIptables version: $(iptables --version | cut -d" " -f2)"

############################################
# CHECK PARAMETERS
############################################
cat "/openvpn-$PROTOCOL-$ENCRYPTION/$REGION.ovpn" &> /dev/null
if [[ "$?" != 0 ]]; then printf "/openvpn-$PROTOCOL-$ENCRYPTION/$REGION.ovpn is not accessible\nSleeping for 10 seconds before exit...\n"; sleep 10; exit 1; fi
# TODO more

############################################
# CHECK FOR TUN DEVICE
############################################
while [ "$(cat /dev/net/tun 2>&1 /dev/null)" != "cat: read error: File descriptor in bad state" ];
do printf "\nTUN device is not opened, sleeping for 30 seconds..."; sleep 30; done
printf "\nTUN device is opened"

############################################
# BLOCKING MALICIOUS HOSTNAMES AND IPs WITH UNBOUND
############################################
touch /etc/unbound/blocks-malicious.conf
printf "\nUnbound malicious hostnames blocking is $BLOCK_MALICIOUS"
if [ "$BLOCK_MALICIOUS" = "on" ] && [ ! -f /etc/unbound/blocks-malicious.conf ]; then
    printf "Extracting malicious hostnames archive..."
    tar -xjf /etc/unbound/malicious-hostnames.bz2 -C /etc/unbound/
    printf "DONE\n"
    printf "Extracting malicious IPs archive..."
    tar -xjf /etc/unbound/malicious-ips.bz2 -C /etc/unbound/
    printf "DONE\n"
    printf "Building blocks-malicious.conf for Unbound..."
    while read hostname; do
        echo "local-zone: \""$hostname"\" static" >> /etc/unbound/blocks-malicious.conf
    done < /etc/unbound/malicious-hostnames
    while read ip; do
        echo "private-address: $ip" >> /etc/unbound/blocks-malicious.conf
    done < /etc/unbound/malicious-ips
    printf "$(cat /etc/unbound/malicious-hostnames | wc -l ) malicious hostnames and $(cat /etc/unbound/malicious-ips | wc -l) malicious IP addresses added\n"
    rm -f /etc/unbound/malicious-hostnames* /etc/unbound/malicious-ips*
else
    touch /etc/unbound/blocks-malicious.conf
fi

############################################
# SETTING DNS OVER TLS TO 1.1.1.1 / 1.0.0.1
############################################
printf "\nLaunching Unbound daemon to connect to Cloudflare DNS 1.1.1.1 at its TLS endpoint..."
unbound
status=$?
if [[ "$status" != 0 ]]; then printf "ERROR with status code $status\nSleeping for 10 seconds before exit...\n"; sleep 10; exit $status; fi
printf "DONE"
printf "\nChanging DNS to localhost..."
echo "nameserver 127.0.0.1" > /etc/resolv.conf
status=$?
if [[ "$status" != 0 ]]; then printf "ERROR with status code $status\nSleeping for 10 seconds before exit...\n"; sleep 10; exit $status; fi
echo "options ndots:0" >> /etc/resolv.conf
printf "DONE"

############################################
# ORIGINAL IP FOR HEALTHCHECK
############################################
printf "\nGetting non VPN public IP address..."
export INITIAL_IP=$(wget -qqO- 'https://duckduckgo.com/?q=what+is+my+ip' | grep -ow 'Your IP address is [0-9.]*[0-9]' | grep -ow '[0-9][0-9.]*')
status=$?
if [[ "$status" != 0 ]]; then printf "ERROR with status code $status\nSleeping for 10 seconds before exit...\n"; sleep 10; exit $status; fi
printf "$INITIAL_IP"

############################################
# FIREWALL
############################################
printf "\nSetting firewall for killswitch purposes..."
printf "\n * Detecting local subnet..."
SUBNET=$(ip route show default | tail -n 1 | awk '// {print $1}')
status=$?
if [[ "$status" != 0 ]]; then printf "ERROR with status code $status\nSleeping for 10 seconds before exit...\n"; sleep 10; exit $status; fi
printf "$SUBNET"
printf "\n * Reading parameters to be used for region $REGION, protocol $PROTOCOL and encryption $ENCRYPTION..."
CONNECTIONSTRING=$(grep -i "/openvpn-$PROTOCOL-$ENCRYPTION/$REGION.ovpn" -e 'privateinternetaccess.com')
status=$?
if [[ "$status" != 0 ]]; then printf "ERROR with status code $status\nSleeping for 10 seconds before exit...\n"; sleep 10; exit $status; fi
PORT=$(echo $CONNECTIONSTRING | cut -d' ' -f3)
if [[ "$PORT" == "" ]]; then printf "Port could not be extracted from configuration file\n"; exit 1; fi
PIADOMAIN=$(echo $CONNECTIONSTRING | cut -d' ' -f2)
if [[ "$PIADOMAIN" == "" ]]; then printf "Port could not be extracted from configuration file\n"; exit 1; fi
sed -i '/^remote $PIADOMAIN $PORT/d' "/openvpn-$PROTOCOL-$ENCRYPTION/$REGION.ovpn" && \
printf "\n   * Port: $PORT"
printf "\n   * Domain: $PIADOMAIN"
printf "\n     * Detecting IP addresses corresponding to $PIADOMAIN..."
VPNIPS=$(nslookup $PIADOMAIN localhost | tail -n +5 | grep -o '[0-9]\{1,3\}\.[0-9]\{1,3\}\.[0-9]\{1,3\}\.[0-9]\{1,3\}')
status=$?
if [[ "$status" != 0 ]]; then printf "ERROR with status code $status\nSleeping for 10 seconds before exit...\n"; sleep 10; exit $status; fi
for ip in $VPNIPS
do
    printf "\n        $ip"
done
printf "\n * Adding IP addresses of $PIADOMAIN to /openvpn-$PROTOCOL-$ENCRYPTION/$REGION.ovpn..."
for ip in $VPNIPS
do
    printf "\n     remote $ip $PORT"
    grep "remote $ip $PORT" "/openvpn-$PROTOCOL-$ENCRYPTION/$REGION.ovpn" || echo "remote $ip $PORT" >> "/openvpn-$PROTOCOL-$ENCRYPTION/$REGION.ovpn"
done
printf "\n * Deleting all iptables rules..."
iptables --flush
status=$?
if [[ "$status" != 0 ]]; then printf "ERROR with status code $status\nSleeping for 10 seconds before exit...\n"; sleep 10; exit $status; fi
iptables --delete-chain
status=$?
if [[ "$status" != 0 ]]; then printf "ERROR with status code $status\nSleeping for 10 seconds before exit...\n"; sleep 10; exit $status; fi
iptables -t nat --flush
status=$?
if [[ "$status" != 0 ]]; then printf "ERROR with status code $status\nSleeping for 10 seconds before exit...\n"; sleep 10; exit $status; fi
iptables -t nat --delete-chain
status=$?
if [[ "$status" != 0 ]]; then printf "ERROR with status code $status\nSleeping for 10 seconds before exit...\n"; sleep 10; exit $status; fi
printf "DONE"
printf "\n * Blocking all output traffic..."
iptables -F OUTPUT
status=$?
if [[ "$status" != 0 ]]; then printf "ERROR with status code $status\nSleeping for 10 seconds before exit...\n"; sleep 10; exit $status; fi
iptables -P OUTPUT DROP
status=$?
if [[ "$status" != 0 ]]; then printf "ERROR with status code $status\nSleeping for 10 seconds before exit...\n"; sleep 10; exit $status; fi
printf "DONE"
printf "\n * Adding rules to accept local loopback traffic..."
iptables -A OUTPUT -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT
status=$?
if [[ "$status" != 0 ]]; then printf "ERROR with status code $status\nSleeping for 10 seconds before exit...\n"; sleep 10; exit $status; fi
iptables -A OUTPUT -o lo -j ACCEPT
status=$?
if [[ "$status" != 0 ]]; then printf "ERROR with status code $status\nSleeping for 10 seconds before exit...\n"; sleep 10; exit $status; fi
printf "DONE"
printf "\n * Adding rules to accept traffic of subnet $SUBNET..."
iptables -A OUTPUT -d $SUBNET -j ACCEPT
status=$?
if [[ "$status" != 0 ]]; then printf "ERROR with status code $status\nSleeping for 10 seconds before exit...\n"; sleep 10; exit $status; fi
printf "DONE"
for ip in $VPNIPS
do
    printf "\n * Adding rules to accept traffic with $ip on port $PROTOCOL $PORT..."
    iptables -A OUTPUT -j ACCEPT -d $ip -o eth0 -p $PROTOCOL -m $PROTOCOL --dport $PORT
    status=$?
    if [[ "$status" != 0 ]]; then printf "ERROR with status code $status\nSleeping for 10 seconds before exit...\n"; sleep 10; exit $status; fi
    printf "DONE"
done
printf "\n * Adding rules to accept traffic going through the tun device..."
iptables -A OUTPUT -o tun0 -j ACCEPT
status=$?
if [[ "$status" != 0 ]]; then printf "ERROR with status code $status\nSleeping for 10 seconds before exit...\n"; sleep 10; exit $status; fi
printf "DONE"

############################################
# USER SECURITY
############################################
printf "\nChanging /auth.conf ownership to nonrootuser with read only access..."
chown nonrootuser /auth.conf
chmod 400 /auth.conf
printf "DONE"

############################################
# OPENVPN LAUNCH
############################################
printf "\nStarting OpenVPN using the following parameters:"
printf "\n * Region: $REGION"
printf "\n * Encryption: $ENCRYPTION"
printf "\n * Protocol: $PROTOCOL"
printf "\n * Port: $PORT"
printf "\n * Initial IP address: $(echo "$VPNIPS" | head -n 1)"
printf "\n\n"
cd "/openvpn-$PROTOCOL-$ENCRYPTION"
openvpn --config "$REGION.ovpn" --user nonrootuser --persist-tun --auth-retry nointeract --auth-user-pass /auth.conf --auth-nocache
status=$?
printf "\nOpenVPN exited with status $status\n"
