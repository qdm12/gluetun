# Maintenance

## With caution

- Remove duplicate `/gluetun` directory creation
- Remove firewall shadowsocks input port?
- Remove `script-security` option
- `ncp-ciphers` to `data-ciphers`
- Remove `ncp-disable`

## Uniformization

- Filter servers by protocol for all
- Multiple IPs addresses support for all proviedrs
- UPDATER_PERIOD only update provider in use

## Code

- Use `github.com/qdm12/ddns-updater/pkg/publicip`
- Windows and Darwin development support
- Use `internal/netlink` in firewall and routing packages

## Features

- Pprof server
- Pre-install DNSSEC files so DoT can be activated even before the tunnel is up
- Gluetun entire logs available at control server, maybe in structured format
- Authentication with the control server
- Get announcement from Github file
- Support multiple connections in custom ovpn

## Gluetun V4

- Remove retro environment variables
- Updater servers version reset to 1
- Change models to all have IPs instead of IP
- Remove HTTP server v0
- `PORT` to `OPENVPN_PORT`
- `UNBLOCK` to `DOT_UNBOUND_UNBLOCK`
- `PROTOCOL` to `OPENVPN_PROTOCOL`
- `PORT_FORWARDING`
- Change servers filtering environment variables to plural
- Remove `WIREGUARD_PORT`
- `WIREGUARD_ADDRESS` to `WIREGUARD_ADDRESSES`
- Only use `custom` VPNSP for custom OpenVPN configurations
- `VPNSP` compulsory
- Change `VPNSP` to `VPN_SERVICE_PROVIDER`
- Change `REGION` (etc.) to `SERVER_REGIONS`
- Remove `PUBLICIP_FILE`
- Remove retro-compatibility where OPENVPN_CONFIG != "" implies VPNSP = "custom"
 and set `OPENVPN_CUSTOM_CONFIG` default to `/gluetun/custom.ovpn`
- Split servers.json and compress it
- Use relative paths everywhere instead of absolute
