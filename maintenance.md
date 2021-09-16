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

## Gluetun pre-v4

- Finish HTTP server v1 or v2
- Renamings
  - `UNBLOCK` to `DOT_UNBOUND_UNBLOCK`
  - `PIA_ENCRYPTION` to `PRIVATE_INTERNET_ACCESS_OPENVPN_ENCRYPTION_PRESET`
  - `PORT_FORWARDING` to `PRIVATE_INTERNET_ACCESS_VPN_PORT_FORWARDING`
  - Rename PIA's `REGION` to `COUNTRY`
  - `WIREGUARD_ADDRESS` to `WIREGUARD_ADDRESSES`
  - `VPNSP` to `VPN_SERVICE_PROVIDER`
  - Rename `REGION` (etc.) to `SERVER_REGIONS`
- Split servers.json and compress it

## Gluetun V4

- Remove retro environment variables:
  - `PORT`
  - `UNBLOCK`
  - `PROTOCOL`
  - `PIA_ENCRYPTION`
  - `PORT_FORWARDING`
  - `WIREGUARD_PORT`
  - `REGION` for PIA
  - `WIREGUARD_ADDRESS`
  - `VPNSP`
  - All old location filters such as `REGION`, `COUNTRY`, etc.
- Remove other retro logic
  - `VPNSP`'s `pia = private ...`
  - Remove `OPENVPN_CONFIG` != "" implies `VPNSP` = "custom" AND set `OPENVPN_CUSTOM_CONFIG` default to `/gluetun/custom.ovpn`
- Remove functionalities
  - `SERVER_NUMBER`
  - `PUBLICIP_FILE`
  - `PORT_FORWARDING_STATUS_FILE`
- Updater servers version reset to 1
- Remove HTTP server v0
- Change to compulsory
  - `VPN_SERVICE_PROVIDER`
- Use relative paths everywhere instead of absolute
