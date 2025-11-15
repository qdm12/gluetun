#!/bin/sh

# Description: This script updates the peer-port for the qBittorrent torrent client using its WebUI API.
#
# How to use:
# 1. (Optional) Disable authentication for localhost clients in qBittorrent WebUI settings ("Bypass authentication for clients on localhost" or `bypass_local_auth` in json).
# 2. Set the environment variable:
# VPN_PORT_FORWARDING_UP_COMMAND=/bin/sh -c "/scripts/qbittorrent-port-update.sh --port {{PORTS}} --webui-port 9081"
# Alternatively, you can use `--user` and `--pass` options to provide credentials.

build_default_url() {
    port="${1:-$WEBUI_PORT}"
    echo "http://127.0.0.1:${port}/api"
}

# default values
WEBUI_PORT="8080"
DEFAULT_URL=$(build_default_url)
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
    echo "  --webui-port PORT  Specify the qBittorrent WebUI Port."
    echo "                     Default: ${WEBUI_PORT}"
    echo "  --url URL          Specify the qBittorrent API URL."
    echo "                     Default: ${DEFAULT_URL}"
    echo "                     Overrides --webui-port option."
    echo "Example:"
    echo "  $0 --user admin --pass **** --port 40409"
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
        PORT=$(echo "$2" | cut -d',' -f1)
        shift 2
        ;;
    --webui-port)
        WEBUI_PORT="$2"
        shift 2
        ;;
    --url)
        PREF_URL="$2"
        shift 2
        ;;
    *)
        echo "Unknown option: $1"
        usage
        exit 1
        ;;
    esac
done

if [ -z "${PORT}" ]; then
    echo "ERROR: --port is required but not provided"
    exit 1
fi

if [ -z "${PREF_URL+x}" ]; then
    PREF_URL=$(build_default_url)
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
        --header "Referer: ${PREF_URL}" \
        --post-data "username=${USERNAME}&password=${PASSWORD}" \
        "${PREF_URL}/v2/auth/login" \
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
wget ${WGET_OPTS} -qO- --post-data="json={\"random_port\":false,\"upnp\":false,\"listen_port\":0}" "$PREF_URL/v2/app/setPreferences"
if [ $? -ne 0 ]; then
    echo "ERROR: Could not update qBittorrent peer-port (first call failed)."
    exit 1
fi

# second call to set the actual port
wget ${WGET_OPTS} -qO- --post-data="json={\"listen_port\":$PORT}" "$PREF_URL/v2/app/setPreferences"
if [ $? -ne 0 ]; then
    echo "ERROR: Could not update qBittorrent peer-port (second call failed)."
    exit 1
fi

echo "Success! qBittorrent peer-port updated to ${PORT}"
exit 0
