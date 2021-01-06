# Gluetun VPN client

*Lightweight swiss-knife-like VPN client to tunnel to Private Internet Access,
Mullvad, Windscribe, Surfshark Cyberghost, VyprVPN, NordVPN, PureVPN and Privado VPN servers, using Go, OpenVPN, iptables, DNS over TLS, ShadowSocks and an HTTP proxy*

**ANNOUNCEMENT**: *New Docker image name `qmcgaw/gluetun`*

<img height="250" src="https://raw.githubusercontent.com/qdm12/gluetun/master/title.svg?sanitize=true">

[![Size](https://img.shields.io/docker/image-size/qmcgaw/gluetun?sort=semver&label=Last%20released%20image)](https://hub.docker.com/r/qmcgaw/gluetun/tags?page=1&ordering=last_updated)
[![Size](https://img.shields.io/docker/image-size/qmcgaw/gluetun/latest?label=Latest%20image)](https://hub.docker.com/r/qmcgaw/gluetun/tags)

[![Docker Pulls](https://img.shields.io/docker/pulls/qmcgaw/private-internet-access.svg)](https://hub.docker.com/r/qmcgaw/private-internet-access)
[![Docker Pulls](https://img.shields.io/docker/pulls/qmcgaw/gluetun.svg)](https://hub.docker.com/r/qmcgaw/gluetun)

![Last release](https://img.shields.io/github/release/qdm12/gluetun?label=Last%20release)
![Last Docker tag](https://img.shields.io/docker/v/qmcgaw/gluetun?sort=semver&label=Last%20Docker%20tag)
![GitHub Release Date](https://img.shields.io/github/release-date/qdm12/gluetun?label=Last%20release%20date)

![Commits since release](https://img.shields.io/github/commits-since/qdm12/gluetun/latest?sort=semver)
[![GitHub last commit](https://img.shields.io/github/last-commit/qdm12/gluetun.svg)](https://github.com/qdm12/gluetun/commits)

[![Lines of code](https://img.shields.io/tokei/lines/github/qdm12/gluetun)](https://github.com/qdm12/gluetun)

## Quick links

- Problem or suggestion?
  - [Start a discussion](https://github.com/qdm12/gluetun/discussions)
  - [Create an issue](https://github.com/qdm12/gluetun/issues)
  - [Check the Wiki](https://github.com/qdm12/gluetun/wiki)
  - [Join the Slack channel](https://join.slack.com/t/qdm12/shared_invite/enQtOTE0NjcxNTM1ODc5LTYyZmVlOTM3MGI4ZWU0YmJkMjUxNmQ4ODQ2OTAwYzMxMTlhY2Q1MWQyOWUyNjc2ODliNjFjMDUxNWNmNzk5MDk)
- Happy?
  - Sponsor me on [github.com/sponsors/qdm12](https://github.com/sponsors/qdm12)
  - Donate to [paypal.me/qmcgaw](https://www.paypal.me/qmcgaw)
  - Drop me [an email](mailto:quentin.mcgaw@gmail.com)

## Features

- Based on Alpine 3.12 for a small Docker image of 52MB
- Supports **Private Internet Access**, **Mullvad**, **Windscribe**, **Surfshark**, **Cyberghost**, **Vyprvpn**, **NordVPN**, **PureVPN** and **Privado** servers
- Supports Openvpn only for now
- DNS over TLS baked in with service provider(s) of your choice
- DNS fine blocking of malicious/ads/surveillance hostnames and IP addresses, with live update every 24 hours
- Choose the vpn network protocol, `udp` or `tcp`
- Built in firewall kill switch to allow traffic only with needed the VPN servers and LAN devices
- Built in Shadowsocks proxy (protocol based on SOCKS5 with an encryption layer, tunnels TCP+UDP)
- Built in HTTP proxy (tunnels HTTP and HTTPS through TCP)
- [Connect other containers to it](https://github.com/qdm12/gluetun/wiki/Connect-to-gluetun)
- [Connect LAN devices to it](https://github.com/qdm12/gluetun/wiki/Connect-to-gluetun)
- Compatible with amd64, i686 (32 bit), **ARM** 64 bit, ARM 32 bit v6 and v7, and even s390x as well as ppc64le 🎆
- VPN server side port forwarding for Private Internet Access and Vyprvpn
- Possibility of split horizon DNS by selecting multiple DNS over TLS providers
- Subprograms all drop root privileges once launched
- Subprograms output streams are all merged together
- Can work as a Kubernetes sidecar container, thanks @rorph

## Setup

1. On some devices you may need to setup your tunnel kernel module on your host with `insmod /lib/modules/tun.ko` or `modprobe tun`
    - [Synology users Wiki page](https://github.com/qdm12/gluetun/wiki/Synology-setup)
1. Launch the container with:

    ```bash
    docker run -d --name gluetun --cap-add=NET_ADMIN \
    -e VPNSP="private internet access" -e REGION="CA Montreal" \
    -e OPENVPN_USER=js89ds7 -e OPENVPN_PASSWORD=8fd9s239G \
    -v /yourpath:/gluetun \
    qmcgaw/gluetun
    ```

    or use [docker-compose.yml](https://github.com/qdm12/gluetun/blob/master/docker-compose.yml) with:

    ```bash
    echo "your openvpn username" > openvpn_user
    echo "your openvpn password" > openvpn_password
    docker-compose up -d
    ```

    You should probably check the many [environment variables](https://github.com/qdm12/gluetun/wiki/Environment-variables) available to adapt the container to your needs.

## Further setup

The following points are all optional but should give you insights on all the possibilities with this container.

- Use [Docker secrets](https://github.com/qdm12/gluetun/wiki/Docker-secrets) to read your credentials instead of environment variables
- [Test your setup](https://github.com/qdm12/gluetun/wiki/Test-your-setup)
- [How to connect other containers and devices to Gluetun](https://github.com/qdm12/gluetun/wiki/Connect-to-gluetun)
- [VPN server side port forwarding](https://github.com/qdm12/gluetun/wiki/Port-forwarding)
- [HTTP control server](https://github.com/qdm12/gluetun/wiki/HTTP-Control-server) to automate things, restart Openvpn etc.
- Update the image with `docker pull qmcgaw/gluetun:latest`. See this [Wiki document](https://github.com/qdm12/gluetun/wiki/Docker-image-tags) for Docker tags available.

## Development

- 💻 [Contribute with code](https://github.com/qdm12/gluetun/wiki/Development) ([existing contributors 👍](https://github.com/qdm12/gluetun/blob/master/.github/CONTRIBUTING.md#Contributors))
- [List of issues and feature requests](https://github.com/qdm12/gluetun/issues)
- [Kanban board](https://github.com/qdm12/gluetun/projects/1)

## License

[![MIT](https://img.shields.io/github/license/qdm12/gluetun)](https://github.com/qdm12/gluetun/master/LICENSE)

## Support

- Sponsor me on [Github](https://github.com/sponsors/qdm12) or donate to [paypal.me/qmcgaw](https://www.paypal.me/qmcgaw)

  [![https://github.com/sponsors/qdm12](https://raw.githubusercontent.com/qdm12/gluetun/master/doc/sponsors.jpg)](https://github.com/sponsors/qdm12)
  [![https://www.paypal.me/qmcgaw](https://raw.githubusercontent.com/qdm12/gluetun/master/doc/paypal.jpg)](https://www.paypal.me/qmcgaw)

- Contribute to the issues and discussions on Github
- Many thanks to @Frepke, @Ralph521, G. Mendez, M. Otmar Weber, J. Perez, A. Cooper and **others** for supporting me financially 🥇👍

## Metadata

[![GitHub commit activity](https://img.shields.io/github/commit-activity/y/qdm12/gluetun.svg)](https://github.com/qdm12/gluetun/commits)
[![GitHub closed PRs](https://img.shields.io/github/issues-pr-closed/qdm12/gluetun.svg)](https://github.com/qdm12/gluetun/pulls?q=is%3Apr+is%3Aclosed)

[![GitHub issues](https://img.shields.io/github/issues/qdm12/gluetun.svg)](https://github.com/qdm12/gluetun/issues)
[![GitHub closed issues](https://img.shields.io/github/issues-closed/qdm12/gluetun.svg)](https://github.com/qdm12/gluetun/issues?q=is%3Aissue+is%3Aclosed)

![Visitors count](https://visitor-badge.laobi.icu/badge?page_id=gluetun.readme)
![GitHub stars](https://img.shields.io/github/stars/qdm12/gluetun?style=social)
![GitHub watchers](https://img.shields.io/github/watchers/qdm12/gluetun?style=social)
![Contributors](https://img.shields.io/github/contributors/qdm12/gluetun)
![GitHub forks](https://img.shields.io/github/forks/qdm12/gluetun?style=social)

![Code size](https://img.shields.io/github/languages/code-size/qdm12/gluetun)
![GitHub repo size](https://img.shields.io/github/repo-size/qdm12/gluetun)

[![dockeri.co](https://dockeri.co/image/qmcgaw/gluetun)](https://hub.docker.com/r/qmcgaw/gluetun)

![Docker Layers for latest](https://img.shields.io/microbadger/layers/qmcgaw/gluetun/latest?label=Docker%20image%20layers)
![Go version](https://img.shields.io/github/go-mod/go-version/qdm12/gluetun)
