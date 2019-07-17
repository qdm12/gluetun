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

printf "[INFO] Reading forwarded port\n"
client_id=`head -n 100 /dev/urandom | sha256sum | tr -d " -"`
exitOnError $? "Unable to generate Client ID"
json=`wget -qO- "http://209.222.18.222:2000/?client_id=$client_id"`
if [ "$json" == "" ]; then
  printf " * Port forwarding is already activated on this connection, has expired, or you are not connected to a PIA region that supports port forwarding\n"
  exit 1
fi
port=`echo $json | jq .port`
port_file="/forwarded_port"
echo "$port" > $port_file
printf " * Written forwarded port to $port_file\n"
ip=`wget -qO- https://diagnostic.opendns.com/myip`
exitOnError $? "Unable to read remote VPN IP"
printf " * Forwarded port is $port on remote VPN IP $ip\n"
printf " * Detecting target VPN interface..."
TARGET_PATH="/openvpn/target"
vpn_device=$(cat $TARGET_PATH/config.ovpn | grep 'dev ' | cut -d" " -f 2)0
exitOnError $? "Unable to find VPN interface"
printf "$vpn_device\n"
printf " * Accepting input traffic through $vpn_device to port $port..."
iptables -A INPUT -i $vpn_device -p tcp --dport $port -j ACCEPT
exitOnError $? "Unable to allow the forwarded port in TCP"
iptables -A INPUT -i $vpn_device -p udp --dport $port -j ACCEPT
exitOnError $? "Unable to allow the forwarded port in UDP"
printf "DONE\n"
