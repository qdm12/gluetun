# Firewall

If you have a strict firewall setup (host level or router level), you might want to setup the following.

## VPN connections

You need the following to allow communicating with the VPN servers

### Private Internet Access

- If `PIA_ENCRYPTION=strong` and `PROTOCOL=udp`: allow outbound UDP 1197 to the corresponding VPN server IPs
- If `PIA_ENCRYPTION=normal` and `PROTOCOL=udp`: allow outbound UDP 1198 to the corresponding VPN server IPs
- If `PIA_ENCRYPTION=strong` and `PROTOCOL=tcp`: allow outbound TCP 501 to the corresponding VPN server IPs
- If `PIA_ENCRYPTION=normal` and `PROTOCOL=tcp`: allow outbound TCP 502 to the corresponding VPN server IPs

### Mullvad

- If `PORT=`, please refer to the mapping of Mullvad servers in [these source code lines](../internal/constants/mullvad.go#L64-L667) to find the corresponding UDP port number and IP address(es) of your choice
- If `PORT=53`, allow outbound UDP 53 to the corresponding VPN server IPs, which you can fine in [the mapping of Mullvad servers](../internal/constants/mullvad.go#L64-L667)
- If `PORT=80`, allow outbound TCP 80 to the corresponding VPN server IPs, which you can fine in [the mapping of Mullvad servers](../internal/constants/mullvad.go#L64-L667)
- If `PORT=443`, allow outbound TCP 443 to the corresponding VPN server IPs, which you can fine in [the mapping of Mullvad servers](../internal/constants/mullvad.go#L64-L667)

### Windscribe

- If `PROTOCOL=udp`: allow outbound UDP 443 to the corresponding VPN server IPs
- If `PROTOCOL=tcp`: allow outbound TCP 1194 to the corresponding VPN server IPs

## Inbound connections

- If `SHADOWSOCKS=on`, allow inbound TCP 8388 and UDP 8388 from your LAN
- If `TINYPROXY=on`, allow inbound TCP 8888 from your LAN
- If you want access to the built-in HTTP control server, allow inbound TCP 8000 from your LAN
