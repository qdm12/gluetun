#!/bin/sh

if ( [ ! -c /dev/net/tun ] ); then
    if ( [ ! -d /dev/net ] ); then
        mkdir -m 755 /dev/net
    fi
    mknod /dev/net/tun c 10 200
fi
if ( !(lsmod | grep -q "^tun\s") ); then
    insmod /lib/modules/tun.ko
fi
