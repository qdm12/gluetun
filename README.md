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
| 19.8MB | 14MB to 80MB | Low to Medium |

<details><summary>Click to show base components</summary><p>

- [Alpine 3.10](https://alpinelinux.org) for a tiny image
- [OpenVPN 2.4.7](https://pkgs.alpinelinux.org/package/v3.10/main/x86_64/openvpn) to tunnel to PIA servers
- [IPtables 1.8.3](https://pkgs.alpinelinux.org/package/v3.10/main/x86_64/iptables) enforces the container to communicate only through the VPN or with other containers in its virtual network (acts as a killswitch)
- [Unbound 1.9.1](https://pkgs.alpinelinux.org/package/v3.10/main/x86_64/unbound) configured with Cloudflare's [1.1.1.1](https://1.1.1.1) DNS over TLS
- [Files and blocking lists built periodically](https://github.com/qdm12/updated/tree/master/files) used with Unbound (see `BLOCK_MALICIOUS` and `BLOCK_NSA` environment variables)

</p></details>

## Features

- <details><summary>Configure everything with environment variables</summary><p>

    - [Destination region](https://www.privateinternetaccess.com/pages/network)
    - Internet protocol
    - Level of encryption
    - Username and password
    - Malicious DNS blocking
    - Extra subnets allowed by firewall
    - Run openvpn without root (but will give reconnect problems)

    </p></details>
- Connect other containers to it, [see this](https://github.com/qdm12/private-internet-access-docker#connect-to-it)
- The *iptables* firewall allows traffic only with needed PIA servers (IP addresses, port, protocol) combinations
- OpenVPN reconnects automatically on failure
- Docker healthcheck pings the DNS 1.1.1.1 to verify the connection is up
- Unbound DNS runs *without root*
- OpenVPN can run *without root* but this disallows OpenVPN reconnecting, it can be set with `NONROOT=yes`
- **ARM** compatible
- Port forwarding

## Setup

1. <details><summary>Requirements</summary><p>

    - A Private Internet Access **username** and **password** - [Sign up](https://www.privateinternetaccess.com/pages/buy-vpn/)
    - Firewall requirements
        - Allow outbound TCP 853 to 1.1.1.1 to allow Unbound to resolve the PIA domain name at start. You can then block it once the container is started.
        - For UDP strong encryption, allow outbound UDP 1197
        - For UDP normal encryption, allow outbound UDP 1198
        - For TCP strong encryption, allow outbound TCP 501
        - For TCP normal encryption, allow outbound TCP 502

    </p></details>

1. Ensure `/dev/net/tun` is setup on your host with either:

    ```sh
    insmod /lib/modules/tun.ko
    # or...
    modprobe tun
    ```

1. <details><summary>CLICK IF YOU HAVE AN ARM DEVICE</summary><p>

    - If you have a ARM 32 bit v6 architecture

        ```sh
        docker build -t qmcgaw/private-internet-access \
        --build-arg BASE_IMAGE=arm32v6/alpine \
        https://github.com/qdm12/private-internet-access-docker.git
        ```

    - If you have a ARM 32 bit v7 architecture

        ```sh
        docker build -t qmcgaw/private-internet-access \
        --build-arg BASE_IMAGE=arm32v7/alpine \
        https://github.com/qdm12/private-internet-access-docker.git
        ```

    - If you have a ARM 64 bit v8 architecture

        ```sh
        docker build -t qmcgaw/private-internet-access \
        --build-arg BASE_IMAGE=arm64v8/alpine \
        https://github.com/qdm12/private-internet-access-docker.git
        ```

    </p></details>

1. Launch the container with:

    ```bash
    docker run -d --name=pia --cap-add=NET_ADMIN --device=/dev/net/tun \
    -e REGION="CA Montreal" -e USER=js89ds7 -e PASSWORD=8fd9s239G \
    qmcgaw/private-internet-access
    ```

    or use [docker-compose.yml](https://github.com/qdm12/private-internet-access-docker/blob/master/docker-compose.yml) with:

    ```bash
    docker-compose up -d
    ```

    Note that you can change all the [environment variables](#environment-variables)

## Testing

Check the PIA IP address matches your expectations

```sh
docker run --rm --network=container:pia alpine:3.10 wget -qO- https://ipinfo.io
```

## Environment variables

| Environment variable | Default | Description |
| --- | --- | --- |
| `REGION` | `CA Montreal` | One of the [PIA regions](https://www.privateinternetaccess.com/pages/network/) |
| `PROTOCOL` | `udp` | `tcp` or `udp` |
| `ENCRYPTION` | `strong` | `normal` or `strong` |
| `USER` | | Your PIA username |
| `PASSWORD` | | Your PIA password |
| `NONROOT` | `no` | Run OpenVPN without root, `yes` or `no` |
| `EXTRA_SUBNETS` | | comma separated subnets allowed in the container firewall (i.e. `192.168.1.0/24,192.168.10.121,10.0.0.5/28`) |
| `BLOCK_MALICIOUS` | `off` | `on` or `off`, blocks malicious hostnames and IPs |
| `BLOCK_NSA` | `off` | `on` or `off`, blocks NSA hostnames |
| `UNBLOCK` | | comma separated string (i.e. `web.com,web2.ca`) to unblock hostnames |

## Connect to it

There are various ways to achieve this, depending on your use case.

- <details><summary>Connect other containers to PIA</summary><p>

    Add `--network=container:pia` when launching the container

    </p></details>
- <details><summary>Connect containers from another docker-compose.yml</summary><p>

    Add `network_mode: "container:pia"` to your *docker-compose.yml*

    </p></details>
- <details><summary>Connect containers in the same docker-compose.yml as PIA</summary><p>

    Add `network_mode: "service:pia"` to your *docker-compose.yml* (no need for `depends_on`)

    </p></details>
- <details><summary>Access ports of containers connected to PIA</summary><p>

    To access port `8000` of container `xyz` and `9000` of container `abc` connected to PIA, you will need a reverse proxy such as `qmcgaw/caddy-scratch` (you can build it for **ARM**, see its [readme](https://github.com/qdm12/caddy-scratch))

    1. Create the file *Caddyfile*

        ```sh
        touch Caddyfile
        chown 1000 Caddyfile
        # chown 1000 because caddy-scratch runs as user ID 1000 by default
        chmod 600 Caddyfile
        ```

        with this content:

        ```ruby
        :8000 {
            proxy / xyz:8000
        }
        :9000 {
            proxy / abc:9000
        }
        ```

        You can of course make more complicated Caddyfile (such as proxying `/xyz` to xyz:8000 and `/abc` to abc:9000, just ask me!)

    1. Run Caddy with

        ```sh
        docker run -d -p 8000:8000/tcp -p 9000:9000/tcp \
        --link pia:xyz --link pia:abc \
        -v $(pwd)/Caddyfile:/Caddyfile:ro \
        qmcgaw/caddy-scratch
        ```

        **WARNING**: Make sure the Docker network in which Caddy runs is the same as the one of PIA. It can be the default `bridge` network.

    1. You can now access xyz:8000 at [localhost:8000](http://localhost:8000) and abc:9000 at [localhost:9000](http://localhost:9000)

    For more containers, add more `--link pia:xxx` and modify the *Caddyfile* accordingly

    If you want to user a *docker-compose.yml*, you can use this example - **make sure PIA is launched and connected first**:

    ```yml
    version: '3'
    services:
      piaproxy:
        image: qmcgaw/caddy-scratch
        container_name: piaproxy
        ports:
          - 8000:8000/tcp
          - 9000:9000/tcp
        external_links:
          - pia:xyz
          - pia:abc
        volumes:
          - ./Caddyfile:/Caddyfile:ro
      abc:
        image: abc
        container_name: abc
        network_mode: "container:pia"
      xyz:
        image: xyz
        container_name: xyz
        network_mode: "container:pia"
    ```

    </p></details>
- <details><summary>Access ports of containers connected to PIA, all in the same docker-compose.yml</summary><p>

    To access port `8000` of container `xyz` and `9000` of container `abc` connected to PIA, you could use:

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
        ports:
          - 8000:8000/tcp
          - 9000:9000/tcp
      abc:
        image: abc
        container_name: abc
        network_mode: "service:pia"
      xyz:
        image: xyz
        container_name: xyz
        network_mode: "service:pia"
    ```

    </p></details>

- <details><summary>Access ports of containers connected to PIA, all in the same docker-compose.yml, using a reverse proxy</summary><p>

    To access port `8000` of container `xyz` and `9000` of container `abc` connected to PIA, you will need a reverse proxy such as `qmcgaw/caddy-scratch` (you can build it for **ARM**, see its [readme](https://github.com/qdm12/caddy-scratch))

    1. Create the file *Caddyfile*

        ```sh
        touch Caddyfile
        chown 1000 Caddyfile
        # chown 1000 because caddy-scratch runs as user ID 1000 by default
        chmod 600 Caddyfile
        ```

        with this content:

        ```ruby
        :8000 {
            proxy / xyz:8000
        }
        :9000 {
            proxy / abc:9000
        }
        ```

        You can of course make more complicated Caddyfile (such as proxying `/xyz` to xyz:8000 and `/abc` to abc:9000, just ask me!)

    1. Use this example:

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
          piaproxy:
            image: qmcgaw/caddy-scratch
            container_name: piaproxy
            ports:
              - 8000:8000/tcp
              - 9000:9000/tcp
            external_links:
              - pia:xyz
              - pia:abc
            volumes:
              - ./Caddyfile:/Caddyfile:ro
          abc:
            image: abc
            container_name: abc
            network_mode: "service:pia"
          xyz:
            image: xyz
            container_name: xyz
            network_mode: "service:pia"
        ```

    </p></details>
- <details><summary>Connect to the PIA through an HTTP proxy (i.e. with Firefox)</summary><p>

    *This is in progress, using Tiny Proxy, thanks for waiting !*

    </p></details>

## Port forwarding

On a running PIA container, say `pia`, simply run:

```sh
docker exec -it pia /portforward.sh
```

And it will indicate you the port forwarded for your current public IP address.

Note that not all regions support port forwarding.

## For the paranoids

- You can review the code which essential consists in the [Dockerfile](https://github.com/qdm12/private-internet-access-docker/blob/master/Dockerfile) and [entrypoint.sh](https://github.com/qdm12/private-internet-access-docker/blob/master/entrypoint.sh)
- Build the images yourself:

    ```bash
    docker build -t qmcgaw/private-internet-access https://github.com/qdm12/private-internet-access-docker.git
    ```

- The download and unziping of PIA openvpn files is done at build for the ones not able to download the zip files
- Checksums for PIA openvpn zip files are not used as these files change often (but HTTPS is used)
- Use `-e ENCRYPTION=strong -e BLOCK_MALICIOUS=on`
- DNS Leaks tests might not work because of [this](https://github.com/qdm12/cloudflare-dns-server#verify-dns-connection) (*TLDR*: DNS server is a local caching intermediary)

## TODOs

- [ ] Tiny proxy for LAN devices to use the container

## License

This repository is under an [MIT license](https://github.com/qdm12/private-internet-access-docker/master/license)
