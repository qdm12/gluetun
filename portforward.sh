#!/bin/sh

client_id=`head -n 100 /dev/urandom | sha256sum | tr -d " -"`
json=`wget -qO- "http://209.222.18.222:2000/?client_id=$client_id" 2>/dev/null`
if [ "$json" == "" ]; then
    printf "Port forwarding is already activated on this connection, has expired, or you are not connected to a PIA region that supports port forwarding\n"
    exit 1
fi
port=`echo $json | grep -Eo [0-9]{3,5}`
ip=`wget -qO- https://diagnostic.opendns.com/myip`
printf "Forwarded port for IP $ip is: $port\n"
