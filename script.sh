#!/bin/sh

set -e
if [ ! -f "/pia/auth.conf" ]]; then
    echo "File auth.conf was not found, aborting !"
    exit 1
fi
openvpn --config "$REGION".ovpn --auth-user-pass auth.conf