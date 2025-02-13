#!/bin/sh

# Description: This script updates the peer port for the qBittorrent torrent client using its WebUI API.
# Note: For this to work, "Bypass authentication for clients on localhost" should be enabled
# in the WebUI settings (json key bypass_local_auth).

build_default_url() {
    port="${1:-$WEBUI_PORT}"
    echo "http://127.0.0.1:${port}/api/v2/app/setPreferences"
}

# default values
WEBUI_PORT="8080"
DEFAULT_URL=$(build_default_url)
WGET_OPTS=""

usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Update qBittorrent peer-port via API."
    echo ""
    echo "Options:"
    echo "  -h, --help             Show this help message and exit."
    echo "  -P, --port PORT        Specify the Transmission Peer Port."
    echo "                         REQUIRED"
    echo "  -W, --webui-port PORT  Specify the qBittorrent WebUI Port."
    echo "                         Default: ${WEBUI_PORT}"
    echo "  -U, --url URL          Specify the qBittorrent API URL."
    echo "                         DEFAULT: ${DEFAULT_URL}"
    echo "                         Overrides --webui option."
    echo "Example:"
    echo "  $0 --port 40409"
    exit 1
}

while [ $# -gt 0 ]; do
    case "$1" in
    -h | --help)
        usage
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
        ;;
    esac
done

if [ -z "${PORT}" ]; then
    echo "ERROR: No PORT provided!"
    exit 1
fi

if [ -z "${PREF_URL+x}" ]; then
    PREF_URL=$(build_default_url)
fi

# update peer host via API
wget --retry-connrefused -qO- ${WGET_OPTS} --post-data="json={\"random_port\":false}" "$PREF_URL"
wget ${WGET_OPTS} -qO- --post-data="json={\"listen_port\":$PORT}" "$PREF_URL"

# check if wget command succeeded
if [ $? -ne 0 ]; then
    echo "ERROR: Could not update peer port"
    exit 1
fi

echo "Success! Peer port updated to ${PORT}"
exit 0
