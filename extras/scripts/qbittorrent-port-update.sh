#!/bin/sh

# Description: This script updates the peer-port for the qBittorrent torrent client using its WebUI API.
# Note: For this to work, "Bypass authentication for clients on localhost" should be enabled
# in the WebUI settings (json key bypass_local_auth).

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
    echo "  -h, --help             Show this help message and exit."
    echo "  -u, --user USER        Specify the qBittorrent username."
    echo "                         (Omit if not required)"
    echo "  -p, --pass PASS        Specify the qBittorrent password."
    echo "                         (Omit if not required)"
    echo "  -P, --port PORT        Specify the qBittorrent peer-port."
    echo "                         REQUIRED"
    echo "  -W, --webui-port PORT  Specify the qBittorrent WebUI Port."
    echo "                         Default: ${WEBUI_PORT}"
    echo "  -U, --url URL          Specify the qBittorrent API URL."
    echo "                         DEFAULT: ${DEFAULT_URL}"
    echo "                         Overrides --webui-port option."
    echo "Example:"
    echo "  $0 -u admin -p **** --port 40409"
}

while [ $# -gt 0 ]; do
    case "$1" in
    -h | --help)
        usage
        exit 0
        ;;
    -u | --user)
        USERNAME="$2"
        _USECRED=true
        shift 2
        ;;
    -p | --pass)
        PASSWORD="$2"
        _USECRED=true
        shift 2
        ;;
    -P | --port)
        PORTS="$2"
        PORT=$(echo "$PORTS" | cut -d',' -f1)
        shift 2
        ;;
    -W | --webui-port)
        WEBUI_PORT="$2"
        shift 2
        ;;
    -U | --url)
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
    echo "ERROR: No qBittorrent peer-port provided!"
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

# update peer host via API
wget ${WGET_OPTS} -qO- --post-data="json={\"random_port\":false}" "$PREF_URL/v2/app/setPreferences"
wget ${WGET_OPTS} -qO- --post-data="json={\"listen_port\":$PORT}" "$PREF_URL/v2/app/setPreferences"

# check if wget command succeeded
if [ $? -ne 0 ]; then
    echo "ERROR: Could not update qBittorrent peer-port."
    exit 1
fi

echo "Success! qBittorrent peer-port updated to ${PORT}"
exit 0
