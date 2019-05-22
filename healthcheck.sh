#!/bin/sh

out="$(ping -W 3 -c 1 -q -s 8 1.1.1.1)"
[ $? != 0 ] || exit 0
printf "$out"
exit 1