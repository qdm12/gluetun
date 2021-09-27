# Gluetun VPN client

*Lightweight swiss-knife-like VPN client to tunnel to Cyberghost, ExpressVPN, FastestVPN,
HideMyAss, IPVanish, IVPN, Mullvad, NordVPN, Privado, Private Internet Access, PrivateVPN,
ProtonVPN, PureVPN, Surfshark, TorGuard, VPNUnlimited, VyprVPN, WeVPN and Windscribe VPN servers
using Go, OpenVPN or Wireguard, iptables, DNS over TLS, ShadowSocks and an HTTP proxy*

**ANNOUNCEMENT**: Wireguard is now supported for all providers supporting it!

![Title image](https://raw.githubusercontent.com/qdm12/gluetun/master/title.svg)

[![Build status](https://github.com/qdm12/gluetun/actions/workflows/ci.yml/badge.svg)](https://github.com/qdm12/gluetun/actions/workflows/ci.yml)

[![Docker pulls qmcgaw/gluetun](https://img.shields.io/docker/pulls/qmcgaw/gluetun.svg)](https://hub.docker.com/r/qmcgaw/gluetun)
[![Docker pulls qmcgaw/private-internet-access](https://img.shields.io/docker/pulls/qmcgaw/private-internet-access.svg)](https://hub.docker.com/r/qmcgaw/gluetun)

[![Docker stars qmcgaw/gluetun](https://img.shields.io/docker/stars/qmcgaw/gluetun.svg)](https://hub.docker.com/r/qmcgaw/gluetun)
[![Docker stars qmcgaw/private-internet-access](https://img.shields.io/docker/stars/qmcgaw/private-internet-access.svg)](https://hub.docker.com/r/qmcgaw/gluetun)

![Last release](https://img.shields.io/github/release/qdm12/gluetun?label=Last%20release)
![Last Docker tag](https://img.shields.io/docker/v/qmcgaw/gluetun?sort=semver&label=Last%20Docker%20tag)
[![Last release size](https://img.shields.io/docker/image-size/qmcgaw/gluetun?sort=semver&label=Last%20released%20image)](https://hub.docker.com/r/qmcgaw/gluetun/tags?page=1&ordering=last_updated)
![GitHub last release date](https://img.shields.io/github/release-date/qdm12/gluetun?label=Last%20release%20date)
![Commits since release](https://img.shields.io/github/commits-since/qdm12/gluetun/latest?sort=semver)

[![Latest size](https://img.shields.io/docker/image-size/qmcgaw/gluetun/latest?label=Latest%20image)](https://hub.docker.com/r/qmcgaw/gluetun/tags)

[![GitHub last commit](https://img.shields.io/github/last-commit/qdm12/gluetun.svg)](https://github.com/qdm12/gluetun/commits/master)
[![GitHub commit activity](https://img.shields.io/github/commit-activity/y/qdm12/gluetun.svg)](https://github.com/qdm12/gluetun/graphs/contributors)
[![GitHub closed PRs](https://img.shields.io/github/issues-pr-closed/qdm12/gluetun.svg)](https://github.com/qdm12/gluetun/pulls?q=is%3Apr+is%3Aclosed)
[![GitHub issues](https://img.shields.io/github/issues/qdm12/gluetun.svg)](https://github.com/qdm12/gluetun/issues)
[![GitHub closed issues](https://img.shields.io/github/issues-closed/qdm12/gluetun.svg)](https://github.com/qdm12/gluetun/issues?q=is%3Aissue+is%3Aclosed)

[![Lines of code](https://img.shields.io/tokei/lines/github/qdm12/gluetun)](https://github.com/qdm12/gluetun)
![Code size](https://img.shields.io/github/languages/code-size/qdm12/gluetun)
![GitHub repo size](https://img.shields.io/github/repo-size/qdm12/gluetun)
![Go version](https://img.shields.io/github/go-mod/go-version/qdm12/gluetun)

![Visitors count](https://visitor-badge.laobi.icu/badge?page_id=gluetun.readme)

## Quick links

- [Setup](#Setup)
- [Features](#Features)
- Problem?
  - [Check the Wiki](https://github.com/qdm12/gluetun/wiki)
  - [Start a discussion](https://github.com/qdm12/gluetun/discussions)
  - [Fix the Unraid template](https://github.com/qdm12/gluetun/discussions/550)
- Suggestion?
  - [Create an issue](https://github.com/qdm12/gluetun/issues)
  - [Join the Slack channel](https://join.slack.com/t/qdm12/shared_invite/enQtOTE0NjcxNTM1ODc5LTYyZmVlOTM3MGI4ZWU0YmJkMjUxNmQ4ODQ2OTAwYzMxMTlhY2Q1MWQyOWUyNjc2ODliNjFjMDUxNWNmNzk5MDk)
- Happy?
  - Sponsor me on [github.com/sponsors/qdm12](https://github.com/sponsors/qdm12)
  - Donate to [paypal.me/qmcgaw](https://www.paypal.me/qmcgaw)
  - Drop me [an email](mailto:quentin.mcgaw@gmail.com)
- Video:

  [![Video Gif](https://i.imgur.com/CetWunc.gif)](https://youtu.be/0F6I03LQcI4)

- [Substack Console interview](https://console.substack.com/p/console-72)

## Features

- Based on Alpine 3.14 for a small Docker image of 31MB
- Supports: **Cyberghost**, **ExpressVPN**, **FastestVPN**, **HideMyAss**, **IPVanish**, **IVPN**, **Mullvad**, **NordVPN**, **Privado**, **Private Internet Access**, **PrivateVPN**, **ProtonVPN**, **PureVPN**,  **Surfshark**, **TorGuard**, **VPNUnlimited**, **Vyprvpn**, **WeVPN**, **Windscribe** servers
- Supports OpenVPN for all providers listed
- Supports Wireguard
  - For **Mullvad**, **Ivpn** and **Windscribe**
  - For **Torguard**, **VPN Unlimited** and **WeVPN** using [the custom provider](https://github.com/qdm12/gluetun/wiki/Environment-variables#custom)
  - For custom Wireguard configurations using [the custom provider](https://github.com/qdm12/gluetun/wiki/Environment-variables#custom)
  - More in progress, see [#134](https://github.com/qdm12/gluetun/issues/134)
- DNS over TLS baked in with service provider(s) of your choice
- DNS fine blocking of malicious/ads/surveillance hostnames and IP addresses, with live update every 24 hours
- Choose the vpn network protocol, `udp` or `tcp`
- Built in firewall kill switch to allow traffic only with needed the VPN servers and LAN devices
- Built in Shadowsocks proxy (protocol based on SOCKS5 with an encryption layer, tunnels TCP+UDP)
- Built in HTTP proxy (tunnels HTTP and HTTPS through TCP)
- [Connect other containers to it](https://github.com/qdm12/gluetun/wiki/Connect-a-container-to-gluetun)
- [Connect LAN devices to it](https://github.com/qdm12/gluetun/wiki/Connect-a-LAN-device-to-gluetun)
- Compatible with amd64, i686 (32 bit), **ARM** 64 bit, ARM 32 bit v6 and v7, and even ppc64le 🎆
- [Custom VPN server side port forwarding for Private Internet Access](https://github.com/qdm12/gluetun/wiki/Private-internet-access#vpn-server-port-forwarding)
- Possibility of split horizon DNS by selecting multiple DNS over TLS providers
- Unbound subprogram drops root privileges once launched
- Can work as a Kubernetes sidecar container, thanks @rorph

## Setup

🎉 There are now instructions specific to each VPN provider with examples to help you get started as quickly as possible!

Go to the [Wiki](https://github.com/qdm12/gluetun/wiki)!

[🐛 Found a bug in the Wiki?!](https://github.com/qdm12/gluetun/issues/new?assignees=&labels=%F0%9F%93%84+Wiki+issue&template=wiki+issue.md&title=Wiki+issue%3A+)

Here's a docker-compose.yml for the laziest:

```yml
version: "3"
services:
  gluetun:
    image: qmcgaw/gluetun
    cap_add:
      - NET_ADMIN
    ports:
      - 8888:8888/tcp # HTTP proxy
      - 8388:8388/tcp # Shadowsocks
      - 8388:8388/udp # Shadowsocks
    volumes:
      - /yourpath:/gluetun
    environment:
      # See https://github.com/qdm12/gluetun/wiki
      - VPNSP=ivpn
      - VPN_TYPE=openvpn
      # OpenVPN:
      - OPENVPN_USER=
      - OPENVPN_PASSWORD=
      # Wireguard:
      # - WIREGUARD_PRIVATE_KEY=wOEI9rqqbDwnN8/Bpp22sVz48T71vJ4fYmFWujulwUU=
      # - WIREGUARD_ADDRESS=10.64.222.21/32
      # Timezone for accurate log times
      - TZ=
```

## License

[![MIT](https://img.shields.io/github/license/qdm12/gluetun)](https://github.com/qdm12/gluetun/master/LICENSE)
