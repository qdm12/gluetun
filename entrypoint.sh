#!/bin/sh

exitOnError(){
  # $1 must be set to $?
  status=$1
  message=$2
  [ "$message" != "" ] || message="Error!"
  if [ $status != 0 ]; then
    printf "$message (status $status)\n"
    exit $status
  fi
}

exitIfUnset(){
  # $1 is the name of the variable to check - not the variable itself
  var="$(eval echo "\$$1")"
  if [ -z "$var" ]; then
    printf "Environment variable $1 is not set\n"
    exit 1
  fi
}

exitIfNotIn(){
  # $1 is the name of the variable to check - not the variable itself
  # $2 is a string of comma separated possible values
  var="$(eval echo "\$$1")"
  for value in ${2//,/ }
  do
    if [ "$var" = "$value" ]; then
      return 0
    fi
  done
  printf "Environment variable $1 cannot be '$var' and must be one of the following: "
  for value in ${2//,/ }
  do
    printf "$value "
  done
  printf "\n"
  exit 1
}

printf " =========================================\n"
printf " =========================================\n"
printf " ============= PIA CONTAINER =============\n"
printf " =========================================\n"
printf " =========================================\n"
printf " == by github.com/qdm12 - Quentin McGaw ==\n\n"

printf "OpenVPN version: $(openvpn --version | head -n 1 | grep -oE "OpenVPN [0-9\.]* " | cut -d" " -f2)\n"
printf "Unbound version: $(unbound -h | grep "Version" | cut -d" " -f2)\n"
printf "Iptables version: $(iptables --version | cut -d" " -f2)\n"

############################################
# CHECK PARAMETERS
############################################
exitIfUnset USER
exitIfUnset PASSWORD
exitIfNotIn ENCRYPTION "normal,strong"
exitIfNotIn PROTOCOL "tcp,udp"
exitIfNotIn NONROOT "yes,no"
cat "/openvpn/$PROTOCOL-$ENCRYPTION/$REGION.ovpn" &> /dev/null
exitOnError $? "/openvpn/$PROTOCOL-$ENCRYPTION/$REGION.ovpn is not accessible"
for SUBNET in ${EXTRA_SUBNETS//,/ }; do
  if [ $(echo "$SUBNET" | grep -Eo '^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(/([0-2]?[0-9])|([3]?[0-1]))?$') = "" ]; then
    printf "Subnet $SUBNET is not a valid IPv4 subnet of the form 255.255.255.255/31 or 255.255.255.255\n"
    exit 1
  fi
done
exitIfNotIn DOT "on,off"
exitIfNotIn BLOCK_MALICIOUS "on,off"
exitIfNotIn BLOCK_NSA "on,off"
if [ "$DOT" == "off" ]; then
  if [ "$BLOCK_MALICIOUS" == "on" ]; then
    printf "DOT is off so BLOCK_MALICIOUS cannot be on\n"
    exit 1
  elif [ "$BLOCK_NSA" == "on" ]; then
    printf "DOT is off so BLOCK_NSA cannot be on\n"
    exit 1
  fi
fi
exitIfNotIn FIREWALL "on,off"

#####################################################
# Writes to protected file and remove USER, PASSWORD
#####################################################
if [ -f /auth.conf ]; then
  printf "/auth.conf already exists\n"
else
  printf "Writing USER and PASSWORD to protected file /auth.conf..."
  echo "$USER" > /auth.conf
  exitOnError $?
  echo "$PASSWORD" >> /auth.conf
  exitOnError $?
  chown nonrootuser /auth.conf
  exitOnError $?
  chmod 400 /auth.conf
  exitOnError $?
  printf "DONE\n"
  printf "Clearing environment variables USER and PASSWORD..."
  unset -v USER
  unset -v PASSWORD
  printf "DONE\n"
fi

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
if [ "$DOT" == "on" ]; then
  printf "Malicious hostnames and ips blocking is $BLOCK_MALICIOUS\n"
  rm -f /etc/unbound/blocks-malicious.conf
  if [ "$BLOCK_MALICIOUS" = "on" ]; then
    tar -xjf /etc/unbound/blocks-malicious.bz2 -C /etc/unbound/
    printf "$(cat /etc/unbound/blocks-malicious.conf | grep "local-zone" | wc -l ) malicious hostnames and $(cat /etc/unbound/blocks-malicious.conf | grep "private-address" | wc -l) malicious IP addresses blacklisted\n"
  else
    echo "" > /etc/unbound/blocks-malicious.conf
  fi
  if [ "$BLOCK_NSA" = "on" ]; then
    tar -xjf /etc/unbound/blocks-nsa.bz2 -C /etc/unbound/
    printf "$(cat /etc/unbound/blocks-nsa.conf | grep "local-zone" | wc -l ) NSA hostnames blacklisted\n"
    cat /etc/unbound/blocks-nsa.conf >> /etc/unbound/blocks-malicious.conf
    rm /etc/unbound/blocks-nsa.conf
    sort -u -o /etc/unbound/blocks-malicious.conf /etc/unbound/blocks-malicious.conf
  fi
  for hostname in ${UNBLOCK//,/ }
  do
    printf "Unblocking hostname $hostname\n"
    sed -i "/$hostname/d" /etc/unbound/blocks-malicious.conf
  done
fi

############################################
# SETTING DNS OVER TLS TO 1.1.1.1 / 1.0.0.1
############################################
printf "DNS over TLS is $DOT\n"
if [ "$DOT" == "on" ]; then
  printf "Launching Unbound daemon to connect to Cloudflare DNS 1.1.1.1 at its TLS endpoint..."
  unbound
  exitOnError $?
  printf "DONE\n"
  printf "Changing DNS to localhost..."
  echo "nameserver 127.0.0.1" > /etc/resolv.conf
  exitOnError $?
  echo "options ndots:0" >> /etc/resolv.conf
  exitOnError $?
  printf "DONE\n"
fi

############################################
# Reading chosen OpenVPN configuration
############################################
printf "Reading configuration for region $REGION, protocol $PROTOCOL and encryption $ENCRYPTION...\n"
CONNECTIONSTRING=$(grep -i "/openvpn/$PROTOCOL-$ENCRYPTION/$REGION.ovpn" -e 'privateinternetaccess.com')
exitOnError $?
PORT=$(echo $CONNECTIONSTRING | cut -d' ' -f3)
if [ "$PORT" = "" ]; then
  printf "Port not found in /openvpn/$PROTOCOL-$ENCRYPTION/$REGION.ovpn\n"
  exit 1
fi
PIADOMAIN=$(echo $CONNECTIONSTRING | cut -d' ' -f2)
if [ "$PIADOMAIN" = "" ]; then
  printf "Domain not found in /openvpn/$PROTOCOL-$ENCRYPTION/$REGION.ovpn\n"
  exit 1
fi
printf " * Port: $PORT\n"
printf " * Domain: $PIADOMAIN\n"
printf "Detecting IP addresses corresponding to $PIADOMAIN...\n"
VPNIPS=$(nslookup $PIADOMAIN localhost | tail -n +3 | grep -o '[0-9]\{1,3\}\.[0-9]\{1,3\}\.[0-9]\{1,3\}\.[0-9]\{1,3\}')
exitOnError $?
for ip in $VPNIPS; do
  printf "   $ip\n";
done

############################################
# Writing target OpenVPN files
############################################
TARGET_PATH="/openvpn/target"
printf "Creating target OpenVPN files in $TARGET_PATH..."
rm -rf $TARGET_PATH/*
cd "/openvpn/$PROTOCOL-$ENCRYPTION"
cp -f *.crt "$TARGET_PATH"
exitOnError $? "Cannot copy crt file to $TARGET_PATH"
cp -f *.pem "$TARGET_PATH"
exitOnError $? "Cannot copy pem file to $TARGET_PATH"
cp -f "$REGION.ovpn" "$TARGET_PATH/config.ovpn"
exitOnError $? "Cannot copy $REGION.ovpn file to $TARGET_PATH"
sed -i "/$CONNECTIONSTRING/d" "$TARGET_PATH/config.ovpn"
exitOnError $? "Cannot delete '$CONNECTIONSTRING' from $TARGET_PATH/config.ovpn"
sed -i '/resolv-retry/d' "$TARGET_PATH/config.ovpn"
exitOnError $? "Cannot delete 'resolv-retry' from $TARGET_PATH/config.ovpn"
for ip in $VPNIPS; do
  echo "remote $ip $PORT" >> "$TARGET_PATH/config.ovpn"
  exitOnError $? "Cannot add 'remote $ip $PORT' to $TARGET_PATH/config.ovpn"
done
# Uses the username/password from this file to get the token from PIA
echo "auth-user-pass /auth.conf" >> "$TARGET_PATH/config.ovpn"
exitOnError $? "Cannot add 'auth-user-pass /auth.conf' to $TARGET_PATH/config.ovpn"
# Reconnects automatically on failure
echo "auth-retry nointeract" >> "$TARGET_PATH/config.ovpn"
exitOnError $? "Cannot add 'auth-retry nointeract' to $TARGET_PATH/config.ovpn"
# Prevents auth_failed infinite loops - make it interact? Remove persist-tun? nobind?
echo "pull-filter ignore \"auth-token\"" >> "$TARGET_PATH/config.ovpn"
exitOnError $? "Cannot add 'pull-filter ignore \"auth-token\"' to $TARGET_PATH/config.ovpn"
# Runs openvpn without root, as nonrootuser if specified
if [ "$NONROOT" = "yes" ]; then
  echo "user nonrootuser" >> "$TARGET_PATH/config.ovpn"
  exitOnError $? "Cannot add 'user nonrootuser' to $TARGET_PATH/config.ovpn"
fi
# Note: TUN device re-opening will restart the container due to permissions
printf "DONE\n"

############################################
# FIREWALL
############################################
printf "Firewall is $FIREWALL\n"
if [ "$FIREWALL" == "on" ]; then
  printf "Setting firewall for killswitch purposes...\n"
  printf " * Detecting local subnet..."
  SUBNET=$(ip route show | tail -n 1 | cut -d" " -f 1)
  exitOnError $?
  printf "$SUBNET\n"
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
  printf " * Block output traffic..."
  iptables -F OUTPUT
  exitOnError $?
  iptables -P OUTPUT DROP
  exitOnError $?
  printf "DONE\n"
  printf " * Accept established and related output traffic..."
  iptables -A OUTPUT -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT
  exitOnError $?
  printf "DONE\n"
  printf " * Accept local loopback output traffic..."
  iptables -A OUTPUT -o lo -j ACCEPT
  exitOnError $?
  printf "DONE\n"
  printf " * Accept output traffic with local subnet $SUBNET..."
  iptables -A OUTPUT -d $SUBNET -j ACCEPT
  exitOnError $?
  printf "DONE\n"
  for EXTRASUBNET in ${EXTRA_SUBNETS//,/ }
  do
    printf " * Accept output traffic with extra subnet $EXTRASUBNET..."
    iptables -A OUTPUT -d $EXTRASUBNET -j ACCEPT
    exitOnError $?
    printf "DONE\n"
  done
  for ip in $VPNIPS; do
    printf " * Accept output traffic to $ip on interface eth0, port $PROTOCOL $PORT..."
    iptables -A OUTPUT -j ACCEPT -d $ip -o eth0 -p $PROTOCOL -m $PROTOCOL --dport $PORT
    exitOnError $?
    printf "DONE\n"
  done
  printf " * Accept all output traffic on tun0 interface..."
  iptables -A OUTPUT -o tun0 -j ACCEPT
  exitOnError $?
  printf "DONE\n"
fi

############################################
# OPENVPN LAUNCH
############################################
printf "Starting OpenVPN using the following parameters:\n"
printf " * Region: $REGION\n"
printf " * Encryption: $ENCRYPTION\n"
printf " * Protocol: $PROTOCOL\n"
printf " * Port: $PORT\n"
printf " * Initial VPN IP address: $(echo "$VPNIPS" | head -n 1)\n\n"
cd "$TARGET_PATH"
openvpn --config config.ovpn
status=$?
printf "\n =========================================\n"
printf " OpenVPN exit with status $status\n"
printf " =========================================\n\n"
