# Private Internet Access Client (OpenVPN on Alpine Linux)

Docker VPN client to private internet access servers using [OpenVPN](https://openvpn.net/) and Iptables on Alpine Linux.

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

[![?](https://images.microbadger.com/badges/image/qmcgaw/private-internet-access.svg)](https://microbadger.com/images/qmcgaw/private-internet-access)
[![?](https://images.microbadger.com/badges/version/qmcgaw/private-internet-access.svg)](https://microbadger.com/images/qmcgaw/private-internet-access)

| Download size | Image size | RAM usage | CPU usage |
| --- | --- | --- | --- |
| 5MB | 8.94MB | 11MB | Low |

It is based on:

- [Alpine 3.8](https://alpinelinux.org)
- [OpenVPN 2.4.6-r3](https://pkgs.alpinelinux.org/package/v3.8/main/x86_64/openvpn)
- [IPtables 1.6.2-r0](https://pkgs.alpinelinux.org/package/v3.8/main/x86_64/iptables)
- CA-Certificates for the healthcheck (through HTTPS)

It requires:

- A Private Internet Access **username** and **password** - [Sign up](https://www.privateinternetaccess.com/pages/buy-vpn/)
- [Docker](https://docs.docker.com/install/) installed on the host

The PIA *.ovpn* configuration files are downloaded from [the PIA website](https://www.privateinternetaccess.com/openvpn/openvpn.zip) when the Docker image is built. You can build the image yourself if you are paranoid.

You might also want to use [my Cloudflare DNS over TLS Docker container](https://hub.docker.com/r/qmcgaw/cloudflare-dns-server/) to connect to any PIA server so that:

- Man-in-the-middle (ISP, hacker, government) can't block you from resolving the PIA server domain name
    *For example, `austria.privateinternetaccess.com` maps to `185.216.34.229`*
- Man-in-the-middle (ISP, hacker, government) can't see to which server you connect nor when.
    *As the domain name are sent to 1.1.1.1 over TLS, there is no way to examine what domains you are asking to be resolved*

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

1. Create a file *auth.conf* in `/yourhostpath` (for example), with:
    - On the first line: your PIA username (i.e. `js89ds7`)
    - On the second line: your PIA password (i.e. `8fd9s239G`)

### Option 1: Using Docker only

1. Run the container with (at least change `/yourhostpath` to your actual path):

    ```bash
    docker run -d --name=pia \
    --cap-add=NET_ADMIN --device=/dev/net/tun --network=pianet \
    -v /yourhostpath/auth.conf:/auth.conf:ro \
    -e REGION=Germany -e PROTOCOL=udp -e ENCRYPTION=normal \
    qmcgaw/private-internet-access
    ```

    Note that you can change `REGION`, `PROTOCOL` and `ENCRYPTION`, see the [Environment variables section](#environment-variables) for more.
1. Wait about 5 seconds for it to connect to the PIA server. You can check with:

    ```bash
    docker logs pia
    ```

1. Follow the [**Testing section**](#testing)

### Option 2: Using Docker Compose

1. Download [**docker-compose.yml**](https://github.com/qdm12/private-internet-access-docker/blob/master/docker-compose.yml)
1. Edit it and change at least `yourpath`
1. Run the container as a daemon in the background with:

    ```bash
    docker-compose up -d
    ```

    Note that you can change `REGION`, `PROTOCOL` and `ENCRYPTION`, see the [Environment variables section](#environment-variables) for more.
1. Wait about 5 seconds for it to connect to the PIA server. You can check with:

    ```bash
    docker logs -f pia
    ```

1. Follow the [**Testing section**](#testing)

## Testing

1. Note that you can simply use the HEALTCHECK provided. The container will stop by itself if the VPN IP is the same as your initial public IP address.

Otherwise you can follow these instructions:

1. Check your host IP address with:

    ```bash
    curl -s ifconfig.co
    ```

1. Run the **curl** Docker container using your *pia* container with:

    ```bash
    docker run --rm --network=container:pia byrnedo/alpine-curl -s ifconfig.co
    ```

    If the displayed IP address appears and is different that your host IP address, the PIA client works !

## Environment variables

| Environment variable | Default | Description |
| --- | --- | --- |
| `REGION` | `Switzerland` | Any one of the [regions supported by private internet access](https://www.privateinternetaccess.com/pages/network/) |
| `PROTOCOL` | `tcp` | `tcp` or `udp` |
| `ENCRYPTION` | `strong` | `normal` or `strong` |

If you know what you're doing, you can change the container name (`pia`), the hostname (`piaclient`) and the network name (`pianet`) as well.

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

- You might want to build the image yourself
- The download and unziping is done at build for the ones not able to download the zip files with their ISPs.
- Checksums for PIA openvpn zip files are not used as these files change often
- You should use strong encryption for the environment variable `ENCRYPTION`
