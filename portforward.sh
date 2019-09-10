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

warnOnError(){
  # $1 must be set to $?
  status=$1
  message=$2
  [ "$message" != "" ] || message="Undefined error"
  if [ $status != 0 ]; then
    printf "[WARNING] $message, with status $status)\n"
  fi
}

printf "[INFO] Reading forwarded port\n"
printf " * Generating client ID...\n"
client_id=`head -n 100 /dev/urandom | sha256sum | tr -d " -"`
exitOnError $? "Unable to generate Client ID"
printf " * Obtaining forward port from PIA server...\n"
json=`wget -qO- "http://209.222.18.222:2000/?client_id=$client_id"`
exitOnError $? "Could not obtain response from PIA server (does your PIA server support port forwarding?)"
if [ "$json" == "" ]; then
  printf "[ERROR] Port forwarding is already activated on this connection, has expired, or you are not connected to a PIA region that supports port forwarding\n"
  exit 1
fi
printf " * Parsing JSON response...\n"
port=`echo $json | jq .port`
exitOnError $? "Cannot find port in JSON response"
printf " * Writing forwarded port to file...\n"
port_status_folder=`dirname "${PORT_FORWARDING_STATUS_FILE}"`
warnOnError $? "Cannot find parent directory of ${PORT_FORWARDING_STATUS_FILE}"
mkdir -p "${port_status_folder}"
warnOnError $? "Cannot create containing directory ${port_status_folder}"
echo "$port" > "${PORT_FORWARDING_STATUS_FILE}"
warnOnError $? "Cannot write port to ${PORT_FORWARDING_STATUS_FILE}"
printf " * Detecting current VPN IP address...\n"
ip=`wget -qO- https://duckduckgo.com/\?q=ip | grep -oE "\b([0-9]{1,3}\.){3}[0-9]{1,3}\b"`
warnOnError $? "Cannot detect remote VPN IP on https://duckduckgo.com"
printf " * Forwarded port accessible at $ip:$port\n"
printf " * Detecting target VPN interface...\n"
vpn_device=$(cat /openvpn/target/config.ovpn | grep 'dev ' | cut -d" " -f 2)0
exitOnError $? "Unable to find VPN interface in /openvpn/target/config.ovpn"
printf " * Accepting input traffic through $vpn_device to port $port...\n"
iptables -A INPUT -i $vpn_device -p tcp --dport $port -j ACCEPT
exitOnError $? "Unable to allow the forwarded port in TCP"
iptables -A INPUT -i $vpn_device -p udp --dport $port -j ACCEPT
exitOnError $? "Unable to allow the forwarded port in UDP"
printf "[INFO] Port forwarded successfully\n"
