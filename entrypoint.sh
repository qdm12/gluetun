#!/bin/sh

exitOnError(){
    # $1 should be $?
    status=$1
    message=$2
    [ "$message" != "" ] || message="Error!"
    if [ $status != 0 ]; then
      printf "$message (status $status)\n"
      exit $status
    fi
}

printf "\n =========================================\n"
printf "=========================================\n"
printf "============= PIA CONTAINER =============\n"
printf "=========================================\n"
printf "=========================================\n"
printf "== by github.com/qdm12 - Quentin McGaw ==\n\n"

printf "OpenVPN version: $(openvpn --version | head -n 1 | grep -oE "OpenVPN [0-9\.]* " | cut -d" " -f2)\n"
printf "Unbound version: $(unbound -h | grep "Version" | cut -d" " -f2)\n"
printf "Iptables version: $(iptables --version | cut -d" " -f2)\n"

############################################
# CHECK PARAMETERS
############################################
cat "/openvpn-$PROTOCOL-$ENCRYPTION/$REGION.ovpn" &> /dev/null
exitOnError $? "/openvpn-$PROTOCOL-$ENCRYPTION/$REGION.ovpn is not accessible"
# TODO more

############################################
# CHECK FOR TUN DEVICE
############################################
while [ "$(cat /dev/net/tun 2>&1 /dev/null)" != "cat: read error: File descriptor in bad state" ]; do
  printf "TUN device is not available, sleeping for 30 seconds...\n"
  sleep 30
done
printf "TUN device OK\n"

############################################
# BLOCKING MALICIOUS HOSTNAMES AND IPs WITH UNBOUND
############################################
touch /etc/unbound/blocks-malicious.conf
printf "Unbound malicious hostnames blocking is $BLOCK_MALICIOUS\n"
if [ "$BLOCK_MALICIOUS" = "on" ] && [ ! -f /etc/unbound/blocks-malicious.conf ]; then
    printf "Extracting malicious hostnames archive..."
    tar -xjf /etc/unbound/malicious-hostnames.bz2 -C /etc/unbound/
    exitOnError $?
    printf "DONE\n"
    printf "Extracting malicious IPs archive..."
    tar -xjf /etc/unbound/malicious-ips.bz2 -C /etc/unbound/
    exitOnError $?
    printf "DONE\n"
    printf "Building blocks-malicious.conf for Unbound..."
    while read hostname; do
        echo "local-zone: \""$hostname"\" static" >> /etc/unbound/blocks-malicious.conf
    done < /etc/unbound/malicious-hostnames
    exitOnError $?
    while read ip; do
        echo "private-address: $ip" >> /etc/unbound/blocks-malicious.conf
    done < /etc/unbound/malicious-ips
    exitOnError $?
    printf "$(cat /etc/unbound/malicious-hostnames | wc -l ) malicious hostnames and $(cat /etc/unbound/malicious-ips | wc -l) malicious IP addresses added\n"
    rm -f /etc/unbound/malicious-hostnames* /etc/unbound/malicious-ips*
else
    touch /etc/unbound/blocks-malicious.conf
fi

############################################
# SETTING DNS OVER TLS TO 1.1.1.1 / 1.0.0.1
############################################
printf "Launching Unbound daemon to connect to Cloudflare DNS 1.1.1.1 at its TLS endpoint...\n"
unbound
exitOnError $?
printf "DONE\n"
printf "Changing DNS to localhost..."
echo "nameserver 127.0.0.1" > /etc/resolv.conf
exitOnError $?
echo "options ndots:0" >> /etc/resolv.conf
exitOnError $?
printf "DONE\n"

############################################
# ORIGINAL IP FOR HEALTHCHECK
############################################
printf "Getting non VPN public IP address..."
export INITIAL_IP=$(wget -qO- 'https://duckduckgo.com/?q=what+is+my+ip' | grep -o 'Your IP address is [0-9.]*[0-9]' | grep -o '[0-9][0-9.]*')
exitOnError $?
printf "$INITIAL_IP\n"

############################################
# FIREWALL
############################################
printf "Setting firewall for killswitch purposes...\n"
printf " * Detecting local subnet..."
SUBNET=$(ip route show default | tail -n 1 | awk '// {print $1}')
exitOnError $?
printf "$SUBNET\n"
printf " * Reading parameters to be used for region $REGION, protocol $PROTOCOL and encryption $ENCRYPTION..."
CONNECTIONSTRING=$(grep -i "/openvpn-$PROTOCOL-$ENCRYPTION/$REGION.ovpn" -e 'privateinternetaccess.com')
exitOnError $?
PORT=$(echo $CONNECTIONSTRING | cut -d' ' -f3)
if [ "$PORT" = "" ]; then
  printf "Port not found in /openvpn-$PROTOCOL-$ENCRYPTION/$REGION.ovpn\n"
  exit 1
fi
PIADOMAIN=$(echo $CONNECTIONSTRING | cut -d' ' -f2)
if [ "$PIADOMAIN" = "" ]; then
  printf "Domain not found in /openvpn-$PROTOCOL-$ENCRYPTION/$REGION.ovpn\n"
  exit 1
fi
sed -i '/^remote $PIADOMAIN $PORT/d' "/openvpn-$PROTOCOL-$ENCRYPTION/$REGION.ovpn"
exitOnError $? "Can't delete remote connection string in /openvpn-$PROTOCOL-$ENCRYPTION/$REGION.ovpn"
printf "DONE\n"
printf "   * Port: $PORT\n"
printf "   * Domain: $PIADOMAIN\n"
printf "     * Detecting IP addresses corresponding to $PIADOMAIN..."
VPNIPS=$(nslookup $PIADOMAIN localhost | tail -n +5 | grep -o '[0-9]\{1,3\}\.[0-9]\{1,3\}\.[0-9]\{1,3\}\.[0-9]\{1,3\}')
exitOnError $?
printf "DONE\n"
for ip in $VPNIPS; do printf "        $ip\n"; done
printf " * Adding IP addresses of $PIADOMAIN to /openvpn-$PROTOCOL-$ENCRYPTION/$REGION.ovpn...\n"
for ip in $VPNIPS; do
  if [ $(grep "remote $ip $PORT" "/openvpn-$PROTOCOL-$ENCRYPTION/$REGION.ovpn") != "" ]; then
    printf "     remote $ip $PORT (already present)\n"
  else
    printf "     remote $ip $PORT\n"
    echo "remote $ip $PORT" >> "/openvpn-$PROTOCOL-$ENCRYPTION/$REGION.ovpn"
  fi
done
printf " * Deleting all iptables rules..."
iptables --flush
exitOnError $?
iptables --delete-chain
exitOnError $?
iptables -t nat --flush
exitOnError $?
iptables -t nat --delete-chain
exitOnError $?
printf "DONE\n"
printf " * Blocking all output traffic..."
iptables -F OUTPUT
exitOnError $?
iptables -P OUTPUT DROP
exitOnError $?
printf "DONE\n"
printf " * Adding rules to accept local loopback traffic..."
iptables -A OUTPUT -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT
exitOnError $?
iptables -A OUTPUT -o lo -j ACCEPT
exitOnError $?
printf "DONE\n"
printf " * Adding rules to accept traffic of subnet $SUBNET..."
iptables -A OUTPUT -d $SUBNET -j ACCEPT
exitOnError $?
printf "DONE\n"
for ip in $VPNIPS; do
    printf " * Adding rules to accept traffic with $ip on port $PROTOCOL $PORT..."
    iptables -A OUTPUT -j ACCEPT -d $ip -o eth0 -p $PROTOCOL -m $PROTOCOL --dport $PORT
    exitOnError $?
    printf "DONE\n"
done
printf " * Adding rules to accept traffic going through the tun device..."
iptables -A OUTPUT -o tun0 -j ACCEPT
exitOnError $?
printf "DONE\n"

############################################
# USER SECURITY
############################################
printf "Changing /auth.conf ownership to nonrootuser with read only access..."
err=$(chown nonrootuser /auth.conf 2>&1)
if [ "$(echo "$err" | grep "Read-only file system")" = "" ]; then exitOnError $?; fi
err=$(chmod 400 /auth.conf 2>&1)
if [ "$(echo "$err" | grep "Read-only file system")" = "" ]; then exitOnError $?; fi
printf "DONE\n"

############################################
# OPENVPN LAUNCH
############################################
printf "Starting OpenVPN using the following parameters:\n"
printf " * Region: $REGION\n"
printf " * Encryption: $ENCRYPTION\n"
printf " * Protocol: $PROTOCOL\n"
printf " * Port: $PORT\n"
printf " * Initial IP address: $(echo "$VPNIPS" | head -n 1)\n\n"
cd "/openvpn-$PROTOCOL-$ENCRYPTION"
exitOnError $? "Can't access /openvpn-$PROTOCOL-$ENCRYPTION"
openvpn --config "$REGION.ovpn" --user nonrootuser --persist-tun --auth-retry nointeract --auth-user-pass /auth.conf --auth-nocache
printf "\nOpenVPN exited with status $?\n"
