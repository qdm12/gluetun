#!/bin/sh

printf "\n ========================================="
printf "\n ========================================="
printf "\n ============= PIA CONTAINER ============="
printf "\n ========================================="
printf "\n ========================================="
printf "\n == by github.com/qdm12 - Quentin McGaw ==\n"

cd /openvpn-$PROTOCOL-$ENCRYPTION

############################################
# CHECK FOR TUN DEVICE
############################################
while [ "$(cat /dev/net/tun 2>&1 /dev/null)" != "cat: read error: File descriptor in bad state" ];
do
    printf "\nTUN device is not opened, sleeping for 30 seconds..."
    sleep 30
done
printf "\nTUN device is opened"

############################################
# BLOCKING MALICIOUS HOSTS WITH UNBOUND
############################################
touch /etc/unbound/blocks-malicious.conf
printf "\nUnbound malicious hosts blocking is $BLOCK_MALICIOUS"
if [[ "$BLOCK_MALICIOUS" == "on" ]]; then
  printf "\nExtracting blocks-malicious.conf.bz2..."
  tar -xjf /etc/unbound/blocks-malicious.conf.bz2 -C /etc/unbound/
  rm /etc/unbound/blocks-malicious.conf.bz2
  printf "DONE"
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
printf "\n   * Port: $PORT"
printf "\n   * Domain: $PIADOMAIN"
printf "\n     * Detecting IP addresses corresponding to $PIADOMAIN..."
VPNIPS=$(nslookup $PIADOMAIN localhost | tail -n +5 | grep -o '[0-9]\{1,3\}\.[0-9]\{1,3\}\.[0-9]\{1,3\}\.[0-9]\{1,3\}')
status=$?
if [[ "$status" != 0 ]]; then printf "ERROR with status code $status\nSleeping for 10 seconds before exit...\n"; sleep 10; exit $status; fi
VPNIPSLENGTH=0
for ip in $VPNIPS
do
    printf "\n        $ip"
    VPNIPSLENGTH=$((VPNIPSLENGTH+1))
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
iptables -F OUTPUT
iptables -P OUTPUT DROP
printf "\n * Adding rules to accept local loopback traffic..."
iptables -A OUTPUT -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT
iptables -A OUTPUT -o lo -j ACCEPT
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
printf "DONE"

############################################
# OPENVPN LAUNCH (retry with next VPN IP if fail)
############################################
failed=1
i=1
PREVIOUSIP=$PIADOMAIN
while [ $failed != 0 ]
do
    VPNIP=$(echo $VPNIPS | cut -d' ' -f$i)
    printf "\nChanging server VPN address $PREVIOUSIP to $VPNIP..."
    sed -i "s/$PREVIOUSIP/$VPNIP/g" $REGION.ovpn
    status=$?
    if [[ "$status" != 0 ]]; then printf "ERROR with status code $status\nSleeping for 10 seconds before exit...\n"; sleep 10; exit $status; fi
    PREVIOUSIP=$VPNIP
    printf "\nStarting OpenVPN using the following parameters:"
    printf "\n * Region: $REGION"
    printf "\n * Encryption: $ENCRYPTION"
    printf "\n * Address: $PROTOCOL://$VPNIP:$PORT"
    printf "\n\n"
    openvpn --config "$REGION.ovpn" --auth-user-pass /auth.conf
    failed=$?
    if [[ $failed != 0 ]]; then
        printf "\n==> Openvpn failed with status code: $failed"
        i=$((i+1))
        if [[ $i -gt $VPNIPSLENGTH ]]; then
            i=0
        fi
    else
        printf "\n==> Openvpn stopped gracefully"
    fi
done
