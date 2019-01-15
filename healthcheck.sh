#!/bin/sh

ping -W 1 -w 1 -q -s 8 1.1.1.1 &> /dev/null
status=$?
if [ $status = 0 ]; then
  exit 0
fi
printf "Pinging 1.1.1.1 resulted in error status code $status"
exit 1