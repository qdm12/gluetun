#!/bin/sh

exitOnError(){
  # $1 must be set to $?
  status=$1
  message=$2
  [ "$message" != "" ] || message="Undefined error"
  if [ $status != 0 ]; then
    printf "[ERROR] $message, with status $status)\n"
    exit $status
  fi
}

client_id=`head -n 100 /dev/urandom | sha256sum | tr -d " -"`
exitOnError $?
json=`wget -qO- "http://209.222.18.222:2000/?client_id=$client_id" 2>/dev/null`
exitOnError $?
if [ "$json" == "" ]; then
    printf "Port forwarding is already activated on this connection, has expired, or you are not connected to a PIA region that supports port forwarding\n"
    exit 1
fi
port=`echo $json | grep -Eo [0-9]{3,5}`
ip=`wget -qO- https://diagnostic.opendns.com/myip`
exitOnError $?
printf "Forwarded port for IP $ip is: $port\n"
printf "Detecting target VPN interface..."
TARGET_PATH="/openvpn/target"
vpn_device=$(cat $TARGET_PATH/config.ovpn | grep 'dev ' | cut -d" " -f 2)0
exitOnError $?
printf "$vpn_device\n"
printf "Accepting input traffic through $vpn_device to port $port..."
iptables -A INPUT -i $vpn_device --dport $PORT -j ACCEPT
exitOnError $?
printf "DONE\n"
