# Private Internet Access Client (OpenVPN+Iptables+DNS over TLS on Alpine Linux)

*VPN client to tunnel to private internet access servers using OpenVPN, IPtables, DNS over TLS and Alpine Linux*

Optionally set the protocol (TCP, UDP) and the level of encryption using Docker environment variables.

A killswitch is implemented with the *iptables* firewall, only allowing traffic with PIA servers on needed ports / protocols.

[![PIA Docker OpenVPN](https://github.com/qdm12/private-internet-access-docker/raw/master/readme/title.png)](https://hub.docker.com/r/qmcgaw/private-internet-access/)

[![Build Status](https://travis-ci.org/qdm12/private-internet-access-docker.svg?branch=master)](https://travis-ci.org/qdm12/private-internet-access-docker)
[![Docker Build Status](https://img.shields.io/docker/build/qmcgaw/private-internet-access.svg)](https://hub.docker.com/r/qmcgaw/private-internet-access)

[![GitHub last commit](https://img.shields.io/github/last-commit/qdm12/private-internet-access-docker.svg)](https://github.com/qdm12/private-internet-access-docker/issues)
[![GitHub commit activity](https://img.shields.io/github/commit-activity/y/qdm12/private-internet-access-docker.svg)](https://github.com/qdm12/private-internet-access-docker/issues)
[![GitHub issues](https://img.shields.io/github/issues/qdm12/private-internet-access-docker.svg)](https://github.com/qdm12/private-internet-access-docker/issues)

[![Docker Pulls](https://img.shields.io/docker/pulls/qmcgaw/private-internet-access.svg)](https://hub.docker.com/r/qmcgaw/private-internet-access)
[![Docker Stars](https://img.shields.io/docker/stars/qmcgaw/private-internet-access.svg)](https://hub.docker.com/r/qmcgaw/private-internet-access)
[![Docker Automated](https://img.shields.io/docker/automated/qmcgaw/private-internet-access.svg)](https://hub.docker.com/r/qmcgaw/private-internet-access)

[![Image size](https://images.microbadger.com/badges/image/qmcgaw/private-internet-access.svg)](https://microbadger.com/images/qmcgaw/private-internet-access)
[![Image version](https://images.microbadger.com/badges/version/qmcgaw/private-internet-access.svg)](https://microbadger.com/images/qmcgaw/private-internet-access)

| Image size | RAM usage | CPU usage |
| --- | --- | --- |
| 15.7MB | 14MB | Low |

It is based on:

- [Alpine 3.8](https://alpinelinux.org) for a tiny image
- [OpenVPN 2.4.6-r3](https://pkgs.alpinelinux.org/package/v3.8/main/x86_64/openvpn) to tunnel to PIA servers
- [IPtables 1.6.2-r0](https://pkgs.alpinelinux.org/package/v3.8/main/x86_64/iptables) enforces the container to communicate only through the VPN or with other containers in its virtual network (killswitch)
- [Unbound 1.7.3-r0](https://pkgs.alpinelinux.org/package/v3.8/main/x86_64/unbound) configured with Cloudflare's [1.1.1.1](https://1.1.1.1) DNS over TLS
- [Malicious hostnames list](https://github.com/qdm12/malicious-hostnames-docker) used with Unbound (see `BLOCK_MALICIOUS` environment variable)
- [Malicious IPs list](https://github.com/qdm12/malicious-ips-docker) used with Unbound (see `BLOCK_MALICIOUS`)

## Extra features

- Connect other containers to it
- Restarts OpenVPN on failure using another IP address corresponding to the PIA server domain name (usually 10 IPs per subdomain name)
- Regular Docker healthchecks using [duckduckgo.com](https://duckduckgo.com) to obtain your current public IP address and compare it with your initial non-VPN IP address
- Openvpn and Unbound do not run as root

## Requirements

- A Private Internet Access **username** and **password** - [Sign up](https://www.privateinternetaccess.com/pages/buy-vpn/)
- [Docker](https://docs.docker.com/install/) installed on the host
- If you use a firewall on the host:
  - Allow outgoing TCP port 853 for Cloudflare DNS over TLS initial resolution of PIA server domain name.
  - Allow outgoing TCP port 443 for querying duckduckgo.com to obtain the initial IP address for the healthcheck.
  - Allow outgoing TCP port 501 for TCP strong encryption
  - Allow outgoing TCP port 502 for TCP normal encryption
  - Allow outgoing UDP port 1197 for UDP strong encryption
  - Allow outgoing UDP port 1198 for UDP normal encryption

## Setup

1. Make sure you have your `/dev/net/tun` device setup on your host with one of the following commands, depending on your OS:

    ```bash
    insmod /lib/modules/tun.ko
    ```

    Or

    ```bash
    sudo modprobe tun
    ```

1. Create a network to be used by this container and other containers connecting to it with:

    ```bash
    docker network create pianet
    ```

1. Create a file *auth.conf* in `./`, with:
    - On the first line: your PIA username (i.e. `js89ds7`)
    - On the second line: your PIA password (i.e. `8fd9s239G`)
1. Launch the container with:

    ```bash
    docker run -d --name=pia -v ./auth.conf:/auth.conf:ro \
    --cap-add=NET_ADMIN --device=/dev/net/tun --network=pianet \
    -e REGION="CA Montreal" -e PROTOCOL=udp -e ENCRYPTION=strong \
    qmcgaw/private-internet-access
    ```

    or use [docker-compose.yml](https://github.com/qdm12/private-internet-access-docker/blob/master/docker-compose.yml) with:

    ```bash
    docker-compose up -d
    ```

    Note that you can change `REGION`, `PROTOCOL` and `ENCRYPTION`, see the [Environment variables section](#environment-variables)
1. Wait about 5 seconds for it to connect to the PIA server. You can check with:

    ```bash
    docker logs -f pia
    ```

1. Follow the [**Testing section**](#testing)

## Testing

You can simply use the Docker healthcheck. The container will mark itself as **unhealthy** if the public IP address is the same as your initial public IP address. Otherwise you can follow these instructions:

1. Check your host IP address with:

    ```bash
    wget -qO- https://ipinfo.io/ip
    ```

1. Run the **curl** Docker container using your *pia* container with:

    ```bash
    docker run --rm --network=container:pia alpine:3.8 wget -qO- https://ipinfo.io/ip
    ```

    If the displayed IP address appears and is different that your host IP address, the PIA client works !

## Environment variables

| Environment variable | Default | Description |
| --- | --- | --- |
| `REGION` | `CA Montreal` | Any one of the [regions supported by private internet access](https://www.privateinternetaccess.com/pages/network/) |
| `PROTOCOL` | `udp` | `tcp` or `udp` |
| `ENCRYPTION` | `strong` | `normal` or `strong` |
| `BLOCK_MALICIOUS` | `off` | `on` or `off` |

If you know what you're doing, you can change the container name (`pia`) and the network name (`pianet`)

## Connect other containers to it

Connect other Docker containers to the PIA VPN connection by adding `--network=container:pia` when launching them.

---

## EXTRA: Access ports of containers connected to the VPN container

You have to use another container acting as a Reverse Proxy such as Nginx.

**Example**:

- We launch a *Deluge* (torrent client) container with name **deluge** connected to the `pia` container with:

    ```bash
    docker run -d --name=deluge --network=container:pia linuxserver/deluge
    ```

- We launch a *Hydra* container with name **hydra** connected to the `pia` container with:

    ```bash
    docker run -d --name=hydra --network=container:pia linuxserver/hydra
    ```

- HTTP User interfaces are accessible at port 8112 for Deluge and 5075 for Hydra

1. Create the Nginx configuration file *nginx.conf*:

    ```txt
    user  nginx;
    worker_processes  1;
    error_log  /var/log/nginx/error.log warn;
    pid        /var/run/nginx.pid;
    events {
        worker_connections  1024;
    }
    http {
        include       /etc/nginx/mime.types;
        default_type  application/octet-stream;
        log_format  main  '$remote_addr - $remote_user [$time_local] "$request" '
                          '$status $body_bytes_sent "$http_referer" '
                          '"$http_user_agent" "$http_x_forwarded_for"';
        access_log  /var/log/nginx/access.log  main;
        sendfile        on;
        keepalive_timeout  65;
        server {
            listen 1001;
            location / {
                proxy_pass http://deluge:8112/;
                proxy_set_header X-Deluge-Base "/";
            }
        }
        server {
            listen 1002;
            location / {
                proxy_pass http://hydra:5075/;
            }
        }
        include /etc/nginx/conf.d/*.conf;
    }
    ```

1. Run the Alpine [Nginx container](https://hub.docker.com/_/nginx) with:

    ```bash
    docker run -d --name=proxypia -p 8001:1001 -p 8002:1002 \
    --network=pianet --link pia:deluge --link pia:hydra \
    -v /mypathto/nginx.conf:/etc/nginx/nginx.conf:ro nginx:alpine
    ```

1. Access the WebUI of Deluge at [localhost:8000](http://localhost:8000)

For more containers, add more `--link pia:xxx` and modify *nginx.conf* accordingly

## EXTRA: For the paranoids

- You might want to build the Docker image yourself
- The download and unziping is done at build for the ones not able to download the zip files through their ISP
- Checksums for PIA openvpn zip files are not used as these files change often
- You should use strong encryption for the environment variable `ENCRYPTION`
- Let me know if you have any extra idea :) !

## TODOs

- [ ] Iptables should change after initial ip address is obtained
- [ ] More checks for environment variables provided
- [ ] Add checks when launching PIA $?
- [ ] VPN server for other devices to go through the tunnel

## License

This repository is under an [MIT license](https://github.com/qdm12/private-internet-access-docker/master/license)
