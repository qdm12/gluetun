#!/bin/sh

printf "Checking Connection\n"
ping -q -c5 privateinternetaccess.com > /dev/null

if [ $? -eq 0 ]
then
	echo "ok"
fi

#failed
exit 1
