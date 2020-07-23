# Private Internet Access Client (qBittorrent+OpenVPN+Iptables+DNS over TLS on alpine/ubuntu)

<p align="center">
  <a href="https://hub.docker.com/r/j4ym0/pia-qbittorrent">
    <img src="https://images.microbadger.com/badges/image/j4ym0/pia-qbittorrent.svg">
  </a>
  <a href="https://github.com/j4ym0/pia-qbittorrent/releases">
    <img alt="latest version" src="https://img.shields.io/github/v/tag/j4ym0/pia-qbittorrent-docker.svg" />
  </a>
  <a href="https://hub.docker.com/r/j4ym0/pia-qbittorrent">
    <img alt="Pulls from DockerHub" src="https://img.shields.io/docker/pulls/j4ym0/pia-qbittorrent.svg?style=flat-square" />
  </a>
</p>


*Lightweight qBittorrent & Private Internet Access VPN client*

[![PIA Docker OpenVPN](https://github.com/j4ym0/pia-qbittorrent-docker/raw/master/readme/title.png)](https://hub.docker.com/r/qmcgaw/private-internet-access/)



<details><summary>Click to show base components</summary><p>

- [Ubuntu 18.04](https://ubuntu.com) for a base image
- [Alpine 3.12.0](https://alpinelinux.org/) for a base image
- [OpenVPN 2.4.4](https://packages.ubuntu.com/bionic/openvpn) to tunnel to PIA servers
- [IPtables 1.6.1](https://packages.ubuntu.com/bionic/iptables) enforces the container to communicate only through the VPN or with other containers in its virtual network (acts as a killswitch)

</p></details>

## Features

- <details><summary>Configure everything with environment variables</summary><p>

    - [Destination region](https://www.privateinternetaccess.com/pages/network)
    - Internet protocol
    - Level of encryption
    - PIA Username and password
    - DNS Servers

    </p></details>
- Self contained qBittorrent
- Exposed webUI
- Downloads & config Volumes
- The *iptables* firewall allows traffic only with needed PIA servers (IP addresses, port, protocol) combinations
- OpenVPN reconnects automatically on failure
- Docker healthcheck pings the PIA DNS 209.222.18.222 and google.com to verify the connection is up


## Setup

1. <details><summary>Requirements</summary><p>

    - A Private Internet Access **username** and **password** - [Sign up referral link](http://www.privateinternetaccess.com/pages/buy-a-vpn/1218buyavpn?invite=U2FsdGVkX1-Ki-3bKiIknvTQB1F-2Tz79e8QkNeh5Zc%2CbPOXkZjc102Clh5ih5-Pa_TYyTU)
    - External firewall requirements, if you have one
        - Allow outbound TCP 853 to 1.1.1.1 to allow Unbound to resolve the PIA domain name at start. You can then block it once the container is started.
        - For UDP strong encryption, allow outbound UDP 1197
        - For UDP normal encryption, allow outbound UDP 1198
        - For TCP strong encryption, allow outbound TCP 501
        - For TCP normal encryption, allow outbound TCP 502
        - For the built-in web HTTP proxy, allow inbound TCP 8888
    - Docker API 1.25 to support `init`

    </p></details>

1. Launch the container with:

    ```bash
    docker run -d --init --name=pia --cap-add=NET_ADMIN -v /My/Downloads/Folder/:/downloads \
    -p 8888:8888 -e REGION="Netherlands" -e USER=xxxxxxx -e PASSWORD=xxxxxxxx \
    j4ym0/pia-qbittorrent
    ```

    Note that you can:
    - Change the many [environment variables](#environment-variables) available
    - Use `-p 8888:8888/tcp` to access the HTTP web proxy 
    - Pass additional arguments to *openvpn* using Docker's command function (commands after the image name)

## Testing

Check the PIA IP address matches your expectations

try [http://checkmyip.torrentprivacy.com/](http://checkmyip.torrentprivacy.com/)

## Environment variables

| Environment variable | Default | Description |
| --- | --- | --- |
| `REGION` | `Netherlands` | One of the [PIA regions](https://www.privateinternetaccess.com/pages/network/) |
| `PROTOCOL` | `udp` | `tcp` or `udp` |
| `ENCRYPTION` | `strong` | `normal` or `strong` |
| `USER` | | Your PIA username |
| `PASSWORD` | | Your PIA password |
| `WEBUI_PORT` | `8888` | `1024` to `65535` internal port for HTTP proxy |
! `DNS_SERVERS` | `209.222.18.222,209.222.18.218` | DNS servers to use, comma separated

## Connect to it

You can connect via your web browser using http://127.0.0.1:8888 or you public ip / LAN if you have forwarding set up

Default username: admin
Default Password: adminadmin

## For the paranoids

- You can review the code which essential consists in the [Dockerfile](https://github.com/j4ym0/pia-qbittorrent-docker/blob/master/Dockerfile) and [entrypoint.sh](https://github.com/j4ym0/pia-qbittorrent-docker/blob/master/entrypoint.sh)
- Any issues please rais them!!
- Build the images yourself:

    ```bash
    docker build -t j4ym0/pia-qbittorrent https://github.com/j4ym0/pia-qbittorrent-docker.git
    ```

- The download and unziping of PIA openvpn files is done at build for the ones not able to download the zip files
- Checksums for PIA openvpn zip files are not used as these files change often (but HTTPS is used)
- Use `-e ENCRYPTION=strong
- DNS Leaks tests seems to be ok, NEED FEEDBACK

## TODOs

- More DNS leack testing
- Edit config from environment vars

## License

This repository is under an [MIT license](https://github.com/j4ym0/pia-qbittorrent-docker/master/license)
