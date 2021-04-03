#!/bin/bash

ip=$(curl -s https://api.ipify.org)
echo "Public IP is $ip"
