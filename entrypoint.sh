#!/bin/sh

exitOnError(){
  # $1 must be set to $?
  status=$1
  message=$2
  [ "$message" != "" ] || message="Undefined error"
  if [ $status != 0 ]; then
    printf "[ERROR] $message, with status $status\n"
    exit $status
  fi
}

exitIfUnset(){
  # $1 is the name of the variable to check - not the variable itself
  var="$(eval echo "\$$1")"
  if [ -z "$var" ]; then
    printf "[ERROR] Environment variable $1 is not set\n"
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
  printf "[ERROR] Environment variable $1 cannot be '$var' and must be one of the following: "
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
printf "TinyProxy version: $(tinyproxy -v | cut -d" " -f2)\n"
printf "ShadowSocks version: $(ss-server --help | head -n 2 | tail -n 1 | cut -d" " -f 2)\n"

############################################
# BACKWARD COMPATIBILITY PARAMETERS
############################################
[ "$PORT_FORWARDING" == "false" ] && PORT_FORWARDING=on
[ "$PORT_FORWARDING" == "true" ] && PORT_FORWARDING=off
if [ -z $TINYPROXY ] && [ ! -z $PROXY ]; then
  TINYPROXY=$PROXY
fi
if [ -z $TINYPROXY_LOG ] && [ ! -z $PROXY_LOG_LEVEL ]; then
  TINYPROXY_LOG=$PROXY_LOG_LEVEL
fi
if [ -z $TINYPROXY_PORT ] && [ ! -z $PROXY_PORT ]; then
  TINYPROXY_PORT=$PROXY_PORT
fi
if [ -z $TINYPROXY_USER ] && [ ! -z $PROXY_USER ]; then
  TINYPROXY_USER=$PROXY_USER
fi
if [ -z $TINYPROXY_PASSWORD ] && [ ! -z $PROXY_PASSWORD ]; then
  TINYPROXY_PASSWORD=$PROXY_PASSWORD
fi

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
for EXTRA_SUBNET in ${EXTRA_SUBNETS//,/ }; do
  if [ $(echo "$EXTRA_SUBNET" | grep -Eo '^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(/([0-2]?[0-9])|([3]?[0-1]))?$') = "" ]; then
    printf "Extra subnet $EXTRA_SUBNET is not a valid IPv4 subnet of the form 255.255.255.255/31 or 255.255.255.255\n"
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
exitIfNotIn PORT_FORWARDING "on,off"
if [ "$PORT_FORWARDING" == "on" ] && [ -z "$PORT_FORWARDING_STATUS_FILE" ]; then
  printf "PORT_FORWARDING is on but PORT_FORWARDING_STATUS_FILE is not set\n"
  exit 1
fi
exitIfNotIn TINYPROXY "on,off"
if [ "$TINYPROXY" == "on" ]; then
  exitIfNotIn TINYPROXY_LOG "Info,Warning,Error,Critical"
  if [ -z $TINYPROXY_PORT ]; then
    TINYPROXY_PORT=8888
  fi
  if [ `echo $TINYPROXY_PORT | grep -E "^[0-9]+$"` != $TINYPROXY_PORT ]; then
    printf "TINYPROXY_PORT is not a valid number\n"
    exit 1
  elif [ $TINYPROXY_PORT -lt 1024 ]; then
    printf "TINYPROXY_PORT cannot be a privileged port under port 1024\n"
    exit 1
  elif [ $TINYPROXY_PORT -gt 65535 ]; then
    printf "TINYPROXY_PORT cannot be a port higher than the maximum port 65535\n"
    exit 1
  fi
  if [ ! -z "$TINYPROXY_USER" ] && [ -z "$TINYPROXY_PASSWORD" ]; then
    printf "TINYPROXY_USER is set but TINYPROXY_PASSWORD is not set\n"
    exit 1
  elif [ -z "$TINYPROXY_USER" ] && [ ! -z "$TINYPROXY_PASSWORD" ]; then
    printf "TINYPROXY_USER is not set but TINYPROXY_PASSWORD is set\n"
    exit 1
  fi
fi
exitIfNotIn SHADOWSOCKS "on,off"
if [ "$SHADOWSOCKS" == "on" ]; then
  exitIfNotIn SHADOWSOCKS "on,off"
  if [ -z $SHADOWSOCKS_PORT ]; then
    SHADOWSOCKS_PORT=8388
  fi
  if [ `echo $SHADOWSOCKS_PORT | grep -E "^[0-9]+$"` != $SHADOWSOCKS_PORT ]; then
    printf "SHADOWSOCKS_PORT is not a valid number\n"
    exit 1
  elif [ $SHADOWSOCKS_PORT -lt 1024 ]; then
    printf "SHADOWSOCKS_PORT cannot be a privileged port under port 1024\n"
    exit 1
  elif [ $SHADOWSOCKS_PORT -gt 65535 ]; then
    printf "SHADOWSOCKS_PORT cannot be a port higher than the maximum port 65535\n"
    exit 1
  fi
  if [ -z $SHADOWSOCKS_PASSWORD ]; then
    printf "SHADOWSOCKS_PASSWORD is not set\n"
    exit 1
  fi
fi

############################################
# SHOW PARAMETERS
############################################
printf "\n"
printf "OpenVPN parameters:\n"
printf " * Region: $REGION\n"
printf " * Encryption: $ENCRYPTION\n"
printf " * Protocol: $PROTOCOL\n"
printf " * Running without root: $NONROOT\n"
printf "DNS over TLS:\n"
printf " * Activated: $DOT\n"
if [ "$DOT" = "on" ]; then
  printf " * Malicious hostnames DNS blocking: $BLOCK_MALICIOUS\n"
  printf " * NSA related DNS blocking: $BLOCK_NSA\n"
  printf " * Unblocked hostnames: $UNBLOCK\n"
fi
printf "Local network parameters:\n"
printf " * Extra subnets: $EXTRA_SUBNETS\n"
printf " * Tinyproxy HTTP proxy: $TINYPROXY\n"
if [ "$TINYPROXY" == "on" ]; then
  printf " * Tinyproxy port: $TINYPROXY_PORT\n"
  tinyproxy_auth=yes
  if [ -z $TINYPROXY_USER ]; then
    tinyproxy_auth=no
  fi
  printf " * Tinyproxy has authentication: $tinyproxy_auth\n"
  unset -v tinyproxy_auth
fi
printf " * ShadowSocks SOCKS5 proxy: $SHADOWSOCKS\n"
printf "PIA parameters:\n"
printf " * Remote port forwarding: $PORT_FORWARDING\n"
[ "$PORT_FORWARDING" == "on" ] && printf " * Remote port forwarding status file: $PORT_FORWARDING_STATUS_FILE\n"
printf "\n"

#####################################################
# Writes to protected file and remove USER, PASSWORD
#####################################################
if [ -f /auth.conf ]; then
  printf "[INFO] /auth.conf already exists\n"
else
  printf "[INFO] Writing USER and PASSWORD to protected file /auth.conf..."
  echo "$USER" > /auth.conf
  exitOnError $?
  echo "$PASSWORD" >> /auth.conf
  exitOnError $?
  chown nonrootuser /auth.conf
  exitOnError $?
  chmod 400 /auth.conf
  exitOnError $?
  printf "DONE\n"
  printf "[INFO] Clearing environment variables USER and PASSWORD..."
  unset -v USER
  unset -v PASSWORD
  printf "DONE\n"
fi

############################################
# CHECK FOR TUN DEVICE
############################################
if [ "$(cat /dev/net/tun 2>&1 /dev/null)" != "cat: read error: File descriptor in bad state" ]; then
  printf "[WARNING] TUN device is not available, creating it..."
  mkdir -p /dev/net
  mknod /dev/net/tun c 10 200
  exitOnError $?
  chmod 0666 /dev/net/tun
  printf "DONE\n"
fi

############################################
# BLOCKING MALICIOUS HOSTNAMES AND IPs WITH UNBOUND
############################################
if [ "$DOT" == "on" ]; then
  rm -f /etc/unbound/blocks-malicious.conf
  if [ "$BLOCK_MALICIOUS" = "on" ]; then
    tar -xjf /etc/unbound/blocks-malicious.bz2 -C /etc/unbound/
    printf "[INFO] $(cat /etc/unbound/blocks-malicious.conf | grep "local-zone" | wc -l ) malicious hostnames and $(cat /etc/unbound/blocks-malicious.conf | grep "private-address" | wc -l) malicious IP addresses blacklisted\n"
  else
    echo "" > /etc/unbound/blocks-malicious.conf
  fi
  if [ "$BLOCK_NSA" = "on" ]; then
    tar -xjf /etc/unbound/blocks-nsa.bz2 -C /etc/unbound/
    printf "[INFO] $(cat /etc/unbound/blocks-nsa.conf | grep "local-zone" | wc -l ) NSA hostnames blacklisted\n"
    cat /etc/unbound/blocks-nsa.conf >> /etc/unbound/blocks-malicious.conf
    rm /etc/unbound/blocks-nsa.conf
    sort -u -o /etc/unbound/blocks-malicious.conf /etc/unbound/blocks-malicious.conf
  fi
  for hostname in ${UNBLOCK//,/ }
  do
    printf "[INFO] Unblocking hostname $hostname\n"
    sed -i "/$hostname/d" /etc/unbound/blocks-malicious.conf
  done
fi

############################################
# SETTING DNS OVER TLS TO 1.1.1.1 / 1.0.0.1
############################################
if [ "$DOT" == "on" ]; then
  printf "[INFO] Launching Unbound to connect to Cloudflare DNS 1.1.1.1 over TLS..."
  unbound
  exitOnError $?
  printf "DONE\n"
  printf "[INFO] Changing DNS to localhost..."
  printf "`sed '/^nameserver /d' /etc/resolv.conf`\nnameserver 127.0.0.1\n" > /etc/resolv.conf
  exitOnError $?
  printf "DONE\n"
fi

############################################
# Reading chosen OpenVPN configuration
############################################
printf "[INFO] Reading OpenVPN configuration...\n"
CONNECTIONSTRING=$(grep -i "/openvpn/$PROTOCOL-$ENCRYPTION/$REGION.ovpn" -e 'privateinternetaccess.com')
exitOnError $?
PORT=$(echo $CONNECTIONSTRING | cut -d' ' -f3)
if [ "$PORT" = "" ]; then
  printf "[ERROR] Port not found in /openvpn/$PROTOCOL-$ENCRYPTION/$REGION.ovpn\n"
  exit 1
fi
PIADOMAIN=$(echo $CONNECTIONSTRING | cut -d' ' -f2)
if [ "$PIADOMAIN" = "" ]; then
  printf "[ERROR] Domain not found in /openvpn/$PROTOCOL-$ENCRYPTION/$REGION.ovpn\n"
  exit 1
fi
printf " * Port: $PORT\n"
printf " * Domain: $PIADOMAIN\n"
printf "[INFO] Detecting IP addresses corresponding to $PIADOMAIN...\n"
VPNIPS=$(nslookup $PIADOMAIN localhost | tail -n +3 | grep -o '[0-9]\{1,3\}\.[0-9]\{1,3\}\.[0-9]\{1,3\}\.[0-9]\{1,3\}')
exitOnError $?
for ip in $VPNIPS; do
  printf "   $ip\n";
done

############################################
# Writing target OpenVPN files
############################################
TARGET_PATH="/openvpn/target"
printf "[INFO] Creating target OpenVPN files in $TARGET_PATH..."
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
echo "mute-replay-warnings" >> "$TARGET_PATH/config.ovpn"
exitOnError $? "Cannot add 'mute-replay-warnings' to $TARGET_PATH/config.ovpn"
# Note: TUN device re-opening will restart the container due to permissions
printf "DONE\n"

############################################
# NETWORKING
############################################
printf "[INFO] Finding network properties...\n"
printf " * Detecting default gateway..."
DEFAULT_GATEWAY=$(ip r | grep 'default via' | cut -d" " -f 3)
exitOnError $?
printf "$DEFAULT_GATEWAY\n"
printf " * Detecting local interface..."
INTERFACE=$(ip r | grep 'default via' | cut -d" " -f 5)
exitOnError $?
printf "$INTERFACE\n"
printf " * Detecting local subnet..."
SUBNET=$(ip r | grep -v 'default via' | grep $INTERFACE | tail -n 1 | cut -d" " -f 1)
exitOnError $?
printf "$SUBNET\n"
for EXTRASUBNET in ${EXTRA_SUBNETS//,/ }
do
  printf " * Adding $EXTRASUBNET as route via $INTERFACE..."
  ip route add $EXTRASUBNET via $DEFAULT_GATEWAY dev $INTERFACE
  exitOnError $?
  printf "DONE\n"
done
printf " * Detecting target VPN interface..."
VPN_DEVICE=$(cat $TARGET_PATH/config.ovpn | grep 'dev ' | cut -d" " -f 2)0
exitOnError $?
printf "$VPN_DEVICE\n"


############################################
# FIREWALL
############################################
printf "[INFO] Setting firewall\n"
printf " * Blocking everyting\n"
printf "   * Deleting all iptables rules..."
iptables --flush
exitOnError $?
iptables --delete-chain
exitOnError $?
iptables -t nat --flush
exitOnError $?
iptables -t nat --delete-chain
exitOnError $?
printf "DONE\n"
printf "   * Block input traffic..."
iptables -P INPUT DROP
exitOnError $?
printf "DONE\n"
printf "   * Block output traffic..."
iptables -F OUTPUT
exitOnError $?
iptables -P OUTPUT DROP
exitOnError $?
printf "DONE\n"
printf "   * Block forward traffic..."
iptables -P FORWARD DROP
exitOnError $?
printf "DONE\n"

printf " * Creating general rules\n"
printf "   * Accept established and related input and output traffic..."
iptables -A OUTPUT -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT
exitOnError $?
iptables -A INPUT -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT
exitOnError $?
printf "DONE\n"
printf "   * Accept local loopback input and output traffic..."
iptables -A OUTPUT -o lo -j ACCEPT
exitOnError $?
iptables -A INPUT -i lo -j ACCEPT
exitOnError $?
printf "DONE\n"

printf " * Creating VPN rules\n"
for ip in $VPNIPS; do
  printf "   * Accept output traffic to VPN server $ip through $INTERFACE, port $PROTOCOL $PORT..."
  iptables -A OUTPUT -d $ip -o $INTERFACE -p $PROTOCOL -m $PROTOCOL --dport $PORT -j ACCEPT
  exitOnError $?
  printf "DONE\n"
done
printf "   * Accept all output traffic through $VPN_DEVICE..."
iptables -A OUTPUT -o $VPN_DEVICE -j ACCEPT
exitOnError $?
printf "DONE\n"

printf " * Creating local subnet rules\n"
printf "   * Accept input and output traffic to and from $SUBNET..."
iptables -A INPUT -s $SUBNET -d $SUBNET -j ACCEPT
iptables -A OUTPUT -s $SUBNET -d $SUBNET -j ACCEPT
printf "DONE\n"
for EXTRASUBNET in ${EXTRA_SUBNETS//,/ }
do
  printf "   * Accept input traffic through $INTERFACE from $EXTRASUBNET to $SUBNET..."
  iptables -A INPUT -i $INTERFACE -s $EXTRASUBNET -d $SUBNET -j ACCEPT
  exitOnError $?
  printf "DONE\n"
  # iptables -A OUTPUT -d $EXTRASUBNET -j ACCEPT
  # iptables -A OUTPUT -o $INTERFACE -s $SUBNET -d $EXTRASUBNET -j ACCEPT
done

############################################
# TINYPROXY LAUNCH
############################################
if [ "$TINYPROXY" == "on" ]; then
  printf "[INFO] Setting TinyProxy log level to $TINYPROXY_LOG..."
  sed -i "/LogLevel /c\LogLevel $TINYPROXY_LOG" /etc/tinyproxy/tinyproxy.conf
  exitOnError $?
  printf "DONE\n"
  printf "[INFO] Setting TinyProxy port to $TINYPROXY_PORT..."
  sed -i "/Port /c\Port $TINYPROXY_PORT" /etc/tinyproxy/tinyproxy.conf
  exitOnError $?
  printf "DONE\n"
  if [ ! -z "$TINYPROXY_USER" ]; then
    printf "[INFO] Setting TinyProxy credentials..."
    echo "BasicAuth $TINYPROXY_USER $TINYPROXY_PASSWORD" >> /etc/tinyproxy/tinyproxy.conf
    unset -v TINYPROXY_USER
    unset -v TINYPROXY_PASSWORD
    printf "DONE\n"
  fi
  tinyproxy -d &
fi

############################################
# SHADOWSOCKS
############################################
if [ "$SHADOWSOCKS" == "on" ]; then
  ARGS="-c /etc/shadowsocks.json"
  if [ "$SHADOWSOCKS_LOG" == " on" ]; then
    printf "[INFO] Setting ShadowSocks logging..."
    ARGS="$ARGS -v"
    printf "DONE\n"
  fi
  printf "[INFO] Setting ShadowSocks port to $SHADOWSOCKS_PORT..."
  jq ".port_password = {\"$SHADOWSOCKS_PORT\":\"\"}" /etc/shadowsocks.json > /tmp/shadowsocks.json && mv /tmp/shadowsocks.json /etc/shadowsocks.json
  exitOnError $?
  printf "DONE\n"
  printf "[INFO] Setting ShadowSocks password..."
  jq ".port_password[\"$SHADOWSOCKS_PORT\"] = \"$SHADOWSOCKS_PASSWORD\"" /etc/shadowsocks.json > /tmp/shadowsocks.json && mv /tmp/shadowsocks.json /etc/shadowsocks.json
  exitOnError $?
  printf "DONE\n"
  ARGS="$ARGS -s `jq --raw-output '.server' /etc/shadowsocks.json`"
  unset -v SERVER
  ARGS="$ARGS -p $SHADOWSOCKS_PORT"
  ARGS="$ARGS -k $SHADOWSOCKS_PASSWORD"
  ss-server $ARGS &
  unset -v ARGS
fi

############################################
# READ FORWARDED PORT
############################################

if [ "$PORT_FORWARDING" == "on" ]; then
  sleep 10 && /portforward.sh &
fi

############################################
# OPENVPN LAUNCH
############################################
printf "[INFO] Launching OpenVPN\n"
cd "$TARGET_PATH"
openvpn --config config.ovpn "$@"
status=$?
printf "\n =========================================\n"
printf " OpenVPN exit with status $status\n"
printf " =========================================\n\n"
