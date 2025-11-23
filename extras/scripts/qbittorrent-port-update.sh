#!/bin/sh

# Description: This script updates the peer-port for the qBittorrent torrent client using its WebUI API.
#
# How to use:
# 1. Provide username and password via `--user` and `--pass` options.
# 2. (Alternative) Disable authentication for localhost clients in qBittorrent WebUI settings ("Bypass authentication for clients on localhost" or `bypass_local_auth` in json).
# 3. Set the environment variable:
# VPN_PORT_FORWARDING_UP_COMMAND=/bin/sh -c "/scripts/qbittorrent-port-update.sh [--user USER --pass PASS] --port {{PORT}} --iface {{VPN_INTERFACE}} --webui-port 9081"

build_default_url() {
    port="${1:-$WEBUI_PORT}"
    echo "http://127.0.0.1:${port}/api"
}

# default values
VPN_PORT=""
VPN_INTERFACE="tun0"
VPN_ADDRESS=""
WEBUI_PORT="8080"
WEBUI_URL=$(build_default_url "$WEBUI_PORT")

# it might take a few tries for qBittorrent to be available (e.g. slow loading with many torrents)
WGET_OPTS="--retry-connrefused --tries=5"

usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Update qBittorrent peer-port via API"
    echo ""
    echo "Options:"
    echo "  --help             Show this help message and exit."
    echo "  --user USER        Specify the qBittorrent username."
    echo "                     (Omit if not required)"
    echo "  --pass PASS        Specify the qBittorrent password."
    echo "                     (Omit if not required)"
    echo "  --port PORT        Specify the qBittorrent peer-port."
    echo "                     REQUIRED"
    echo "  --iface IFACE      Specify the VPN interface to bind to."
    echo "                     Default: \"${VPN_INTERFACE}\""
    echo "  --addr ADDR        Specify the VPN interface address to bind to."
    echo "                     Available options: empty = All addresses, 0.0.0.0 - All IPv4 addresses, :: - All IPv6 addresses, or a specific IP address"
    echo "                     Default: \"${VPN_ADDRESS}\""
    echo "  --webui-port PORT  Specify the qBittorrent WebUI Port. Not compatible with --url."
    echo "                     Default: \"${WEBUI_PORT}\""
    echo "  --url URL          Specify the qBittorrent API URL. Not compatible with --webui-port."
    echo "                     Default: \"${WEBUI_URL}\""
    echo "Example:"
    echo "  $0 --user ADMIN --pass **** --port 40409"
}

while [ $# -gt 0 ]; do
    case "$1" in
    --help)
        usage
        exit 0
        ;;
    --user)
        USERNAME="$2"
        _USECRED=true
        shift 2
        ;;
    --pass)
        PASSWORD="$2"
        _USECRED=true
        shift 2
        ;;
    --port)
        VPN_PORT=$(echo "$2" | cut -d',' -f1)
        shift 2
        ;;
    --iface)
        VPN_INTERFACE="$2"
        shift 2
        ;;
    --addr)
        VPN_ADDRESS="$2"
        shift 2
        ;;
    --webui-port)
        WEBUI_PORT="$2"
        WEBUI_URL=$(build_default_url "$WEBUI_PORT")
        shift 2
        ;;
    --url)
        WEBUI_URL="$2"
        shift 2
        ;;
    *)
        echo "Unknown option: $1"
        usage
        exit 1
        ;;
    esac
done

if [ -z "${VPN_PORT}" ]; then
    echo "ERROR: --port is required but not provided"
    exit 1
fi

if [ "${_USECRED}" ]; then
    # make sure username AND password were provided
    if [ -z "${USERNAME}" ]; then
        echo "ERROR: qBittorrent username not provided."
        exit 1
    fi
    if [ -z "${PASSWORD}" ]; then
        echo "ERROR: qBittorrent password not provided."
        exit 1
    fi

    cookie=$(wget ${WGET_OPTS} -qO- \
        --header "Referer: ${WEBUI_URL}" \
        --post-data "username=${USERNAME}&password=${PASSWORD}" \
        "${WEBUI_URL}/v2/auth/login" \
        --server-response 2>&1 | \
        grep -i "set-cookie:" | \
        sed 's/.*set-cookie: //I;s/;.*//')

    if [ -z "${cookie}" ]; then
        echo "ERROR: Could not authenticate with qBittorrent."
        exit 1
    fi

    # set cookie for future requests
    WGET_OPTS="${WGET_OPTS} --header=Cookie:$cookie"
fi

# update peer host via API, 0 is a dummy port, required due to https://github.com/qdm12/gluetun-wiki/pull/147
wget ${WGET_OPTS} -qO- --post-data="json={\"random_port\":false,\"upnp\":false,\"listen_port\":0,\"current_network_interface\":\"lo\",\"current_interface_address\":\"\"}" "$WEBUI_URL/v2/app/setPreferences"
if [ $? -ne 0 ]; then
    echo "ERROR: Could not update qBittorrent peer-port (first call failed)."
    exit 1
fi

# second call to set the actual port
wget ${WGET_OPTS} -qO- --post-data="json={\"listen_port\":$VPN_PORT,\"current_network_interface\":\"$VPN_INTERFACE\",\"current_interface_address\":\"$VPN_ADDRESS\"}" "$WEBUI_URL/v2/app/setPreferences"
if [ $? -ne 0 ]; then
    echo "ERROR: Could not update qBittorrent peer-port (second call failed)."
    exit 1
fi

echo "Success! qBittorrent peer-port updated to ${VPN_PORT}"
exit 0
