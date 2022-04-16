# Maintenance

- Rename `UNBLOCK` to `DOT_UNBOUND_UNBLOCK`
- Move constants.*RegionChoices to a validation package
- Common filtering functions
- Refactor providers code to have one directory per VPN provider
- Use DNS v2 beta
- Go 1.18
  - gofumpt
  - Use netip
- Split servers.json and compress it, use Git LFS
- DNS block lists as LFS and built in image
- Finish HTTP server v1 or v2
- Use `github.com/qdm12/ddns-updater/pkg/publicip`
- Windows and Darwin development support

## Features

- Authentication with the control server
- Get announcement from Github file
- Support multiple connections in custom ovpn
- Automate IPv6 detection for OpenVPN

## Gluetun V4

- Remove retro environment variables:
  - `PORT`
  - `UNBLOCK`
  - `PROTOCOL`
  - `PIA_ENCRYPTION`
  - `PORT_FORWARDING`
  - `WIREGUARD_PORT`
  - `REGION` for PIA, Cyberghost
  - `WIREGUARD_ADDRESS`
  - `VPNSP`
  - All old location filters such as `REGION`, `COUNTRY`, etc.
- Remove other retro logic
  - `VPNSP`'s `pia = private ...`
  - Remove `OPENVPN_CONFIG` != "" implies `VPNSP` = "custom" AND set `OPENVPN_CUSTOM_CONFIG` default to `/gluetun/custom.ovpn`
- Remove functionalities
  - `SERVER_NUMBER`
  - `SERVER_NAME`
  - `PUBLICIP_FILE`
  - `PORT_FORWARDING_STATUS_FILE`
- Updater servers version reset to 1
- Reset HTTP server version to v1 and remove older ones
- Change to compulsory
  - `VPN_SERVICE_PROVIDER`
- Use relative paths everywhere instead of absolute
