#!/bin/sh

printf "Checking Connection\n"
out="$(ping -W 3 -c 4 -s 8 209.222.18.222)"
printf "$out"
printf "\n"
printf "Checking DNS\n"
out="$(ping -W 3 -c 4 -s 8 google.co.uk)"
printf "$out"
printf "\n"

printf "VPN Status\n"
vpn="$(netstat -i)"
printf "$vpn \n\n"
exit 1