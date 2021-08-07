# Maintenance

## With caution

- Remove duplicate `/gluetun` directory creation
- Remove firewall shadowsocks input port?
- Re-add `persist-tun`? Run openvpn without root?
- Remove `script-security` option

## Uniformization

- Filter servers by protocol for all
- Multiple IPs addresses support for all proviedrs

## Code

- Use `github.com/qdm12/ddns-updater/pkg/publicip`
- Change firewall debug logs to use `logger.Debug` instead of `fmt.Println`

## Features

- Pprof server
- Pre-install DNSSEC files so DoT can be activated even before the tunnel is up
- Gluetun entire logs available at control server, maybe in structured format
- Authentication with the control server

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
