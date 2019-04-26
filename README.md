# Private Internet Access Client (OpenVPN+Iptables+DNS over TLS on Alpine Linux)

*Lightweight VPN client to tunnel to private internet access servers*

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

[![Donate PayPal](https://img.shields.io/badge/Donate-PayPal-green.svg)](https://paypal.me/qdm12)

| Image size | RAM usage | CPU usage |
| --- | --- | --- |
| 19.6MB | 14MB to 80MB | Low to Medium |

<details><summary>Click to show base components</summary><p>

- [Alpine 3.9](https://alpinelinux.org) for a tiny image
- [OpenVPN 2.4.6-r3](https://pkgs.alpinelinux.org/package/v3.9/main/x86_64/openvpn) to tunnel to PIA servers
- [IPtables 1.6.2-r0](https://pkgs.alpinelinux.org/package/v3.9/main/x86_64/iptables) enforces the container to communicate only through the VPN or with other containers in its virtual network (acts as a killswitch)
- [Unbound 1.7.3-r0](https://pkgs.alpinelinux.org/package/v3.9/main/x86_64/unbound) configured with Cloudflare's [1.1.1.1](https://1.1.1.1) DNS over TLS
- [Files and blocking lists built periodically](https://github.com/qdm12/updated/tree/master/files) used with Unbound (see `BLOCK_MALICIOUS` and `BLOCK_NSA` environment variables)

</p></details>

## Extra features

- Configure everything with environment variables
  - [Destination region](https://www.privateinternetaccess.com/pages/network)
  - Internet protocol
  - Level of encryption
  - Username and password
  - Malicious DNS blocking
  - Extra subnets allowed by firewall
  - Run openvpn without root (but will give reconnect problems)
- Connect other containers to it
- The *iptables* firewall allows traffic only with needed PIA servers (IP addresses, port, protocol) combination
- OpenVPN restarts on failure using another PIA IP address for the same region
- Docker healthcheck uses [https://diagnostic.opendns.com/myip](https://diagnostic.opendns.com/myip) to check that the current public IP address exists in the selected OpenVPN configuration file
- Openvpn and Unbound do not run as root

## Requirements

- A Private Internet Access **username** and **password** - [Sign up](https://www.privateinternetaccess.com/pages/buy-vpn/)
- [Docker](https://docs.docker.com/install/) installed on the host
- If you use a strict firewall on the host/router:
  - Allow outbound TCP 853 to 1.1.1.1 to allow Unbound to resolve the PIA domain name at start. You can then block it once the container is started.
  - For UDP strong encryption, allow outbound UDP 1197
  - For UDP normal encryption, allow outbound UDP 1198
  - For TCP strong encryption, allow outbound TCP 501
  - For TCP normal encryption, allow outbound TCP 502

## Setup

1. Make sure you have your `/dev/net/tun` device setup on your host with one of the following commands, depending on your OS:

    ```bash
    insmod /lib/modules/tun.ko
    ```

    Or

    ```bash
    modprobe tun
    ```

1. **IF YOU HAVE AN ARM DEVICE, depending on your cpu architecture:** replace `qmcgaw/private-internet-access`
   with either `qmcgaw/private-internet-access:armhf` (32 bit) or `qmcgaw/private-internet-access:aarch64` (64 bit).

1. Launch the container with:

    ```bash
    docker run -d --name=pia -v ./auth.conf:/auth.conf:ro \
    --cap-add=NET_ADMIN --device=/dev/net/tun \
    -e REGION="CA Montreal" -e PROTOCOL=udp -e ENCRYPTION=strong \
    -e USER=js89ds7 -e PASSWORD=8fd9s239G \
    qmcgaw/private-internet-access
    ```

    or use [docker-compose.yml](https://github.com/qdm12/private-internet-access-docker/blob/master/docker-compose.yml) with:

    ```bash
    docker-compose up -d
    ```

    Note that you can change all the [environment variables](#environment-variables)
1. Wait about 5 seconds for it to connect to the PIA server. You can check with:

    ```bash
    docker logs -f pia
    ```

1. Follow the [**Testing section**](#testing)

## Testing

You can simply use the Docker healthcheck. The container will mark itself as **unhealthy** if the public IP address is not part of the PIA IPs. Otherwise you can follow these instructions:

1. Check your host IP address with:

    ```bash
    wget -qO- https://ipinfo.io/ip
    ```

1. Run the same command in a Docker container using your *pia* container as network with:

    ```bash
    docker run --rm --network=container:pia alpine:3.9 wget -qO- https://ipinfo.io/ip
    ```

    If the displayed IP address appears and is different that your host IP address, the PIA client works !

## Environment variables

| Environment variable | Default | Description |
| --- | --- | --- |
| `REGION` | `CA Montreal` | One of the [PIA regions](https://www.privateinternetaccess.com/pages/network/) |
| `PROTOCOL` | `udp` | `tcp` or `udp` |
| `ENCRYPTION` | `strong` | `normal` or `strong` |
| `USER` | | Your PIA username |
| `PASSWORD` | | Your PIA password |
| `NONROOT` | `no` | Run OpenVPN without root, `yes` or other |
| `EXTRA_SUBNETS` | | comma separated subnets allowed in the container firewall (i.e. `192.168.1.0/24,192.168.10.121,10.0.0.5/28`) |
| `BLOCK_MALICIOUS` | `off` | `on` or `off`, blocks malicious hostnames and IPs |
| `BLOCK_NSA` | `off` | `on` or `off`, blocks NSA hostnames |
| `UNBLOCK` | | comma separated string (i.e. `web.com,web2.ca`) to unblock hostnames |

## Connect other containers to it

Connect other Docker containers to the PIA VPN connection by adding `--network=container:pia` when launching them.

For containers in the same `docker-compose.yml` as PIA, you can use `network: "service:pia"` (see below)

### Access ports of PIA-connected containers

1. For example, the following containers are launched connected to PIA:

    ```bash
    docker run -d --name=deluge --network=container:pia linuxserver/deluge
    docker run -d --name=hydra --network=container:pia linuxserver/hydra
    ```

    We want to access:
        - The HTTP web UI of Deluge at port **8112**
        - The HTTP Web UI of Hydra at port **5075**

1. In this case we use Nginx for its small size. Create `./nginx.conf` with:

    ```bash
    # nginx.conf
    user nginx;
    worker_processes 1;
    events {
      worker_connections 64;
    }
    http {
      server {
        listen 8000;
        location /deluge {
          proxy_pass http://deluge:8112/;
          proxy_set_header X-Deluge-Base "/deluge";
        }
      }
      server {
        listen 8001;
        location / {
          proxy_pass http://hydra:5075/;
        }
      }
    }
    ```

1. Run the [Nginx Alpine container](https://hub.docker.com/_/nginx):

    ```bash
    docker run -d -p 8000:8000/tcp -p 8001:8001/tcp \
    --link pia:deluge --link pia:hydra \
    -v $(pwd)/nginx.conf:/etc/nginx/nginx.conf:ro \
    nginx:alpine
    ```

    **WARNING**: Make sure the Docker network in which Nginx runs is the same as the one of PIA. It can be the default `bridge` network.

1. Access the WebUI of Deluge at [localhost:8000](http://localhost:8000) and Hydra at [localhost:8001](http://localhost:8001)

For more containers, add more `--link pia:xxx` and modify *nginx.conf* accordingly

The docker compose file would look like:

```yml
version: '3'
services:
  pia:
    image: qmcgaw/private-internet-access
    container_name: pia
    cap_add:
      - NET_ADMIN
    devices:
      - /dev/net/tun
    environment:
      - USER=js89ds7
      - PASSWORD=8fd9s239G
      - PROTOCOL=udp
      - ENCRYPTION=strong
      - REGION=CA Montreal
      - EXTRA_SUBNETS=
      - NONROOT=
    restart: always
  nginx:
    image: nginx:alpine
    container_name: pia_proxy
    ports:
      - 8001:8001/tcp
      - 8002:8002/tcp
    links:
      - pia:deluge
      - pia:hydra
    depends_on:
      - pia
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf:ro
  deluge:
    image: linuxserver/deluge
    container_name: deluge
    network_mode: "container:pia"
    depends_on:
      - pia
    # add more volumes etc.
  hydra:
    image: linuxserver/hydra
    container_name: hydra
    network_mode: "container:hydra"
    depends_on:
      - pia
    # add more volumes etc.
```

## For the paranoids

- You can review the code which essential consits in the [Dockerfile](https://github.com/qdm12/private-internet-access-docker/blob/master/Dockerfile) and [entrypoint.sh](https://github.com/qdm12/private-internet-access-docker/blob/master/entrypoint.sh)
- Build the images yourself:

    ```bash
    docker build -t qmcgaw/private-internet-access https://github.com/qdm12/private-internet-access-docker.git
    ```

- The download and unziping of PIA openvpn files is done at build for the ones not able to download the zip files
- Checksums for PIA openvpn zip files are not used as these files change often (but HTTPS is used)
- Use `-e ENCRYPTION=strong -e BLOCK_MALICIOUS=on`
- DNS Leaks tests might not work because of [this](https://github.com/qdm12/cloudflare-dns-server#verify-dns-connection) (*TLDR*: DNS server is a local caching intermediary)

## TODOs

- [ ] SOCKS/HTTP proxy or VPN server for LAN devices to use the container
- [ ] Travis CI for arm images
- [ ] Nginx scratch
- [ ] Port forwarding

## License

This repository is under an [MIT license](https://github.com/qdm12/private-internet-access-docker/master/license)
