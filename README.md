# Private Internet Access Docker (OpenVPN, Alpine)

Docker VPN client to private internet access servers based on [Alpine Linux](https://alpinelinux.org/) and [OpenVPN](https://openvpn.net/)

[![PIA Docker OpenVPN](https://github.com/qdm12/private-internet-access-docker/raw/master/readme/title.png)](https://hub.docker.com/r/qmcgaw/private-internet-access/)

[![Build Status](https://travis-ci.org/qdm12/private-internet-access-docker.svg?branch=master)](https://travis-ci.org/qdm12/private-internet-access-docker)

[![](https://images.microbadger.com/badges/image/qmcgaw/private-internet-access.svg)](https://microbadger.com/images/qmcgaw/private-internet-access)
[![](https://images.microbadger.com/badges/version/qmcgaw/private-internet-access.svg)](https://microbadger.com/images/qmcgaw/private-internet-access)
| Download size | Image size | RAM usage | CPU usage |
| --- | --- | --- | --- |
| 3.3MB | 8.02MB | 4.3MB | Very low |

It requires:
- A Private Internet Access **username** and **password** - [Sign up](https://www.privateinternetaccess.com/pages/buy-vpn/)
- [Docker](https://docs.docker.com/install/) installed on the host

The PIA configuration files are downloaded from [the PIA website](https://www.privateinternetaccess.com/openvpn/openvpn.zip) when the Docker image gets built.

## Installation & Testing

1. Run the [**tun.sh**](https://raw.githubusercontent.com/qdm12/private-internet-access-docker/master/tun.sh) script on your host machine to ensure you have the **tun** device setup

    ```bash
    wget https://raw.githubusercontent.com/qdm12/private-internet-access-docker/master/tun.sh
    sudo chmod +x tun.sh
    ./tun.sh
    ```

1. Create a file *auth.conf* in `/yourhostpath` (for example), with:
    - On the first line: your PIA username (i.e. `js89ds7`)
    - On the second line: your PIA password (i.e. `8fd9s239G`)
    
### Using Docker only

1. Test the container by connecting another container to it
    1. Run the container interactively with (and change `/yourhostpath/auth.conf`):

        ```bash
        docker run --rm --name=piaTEST --cap-add=NET_ADMIN \
        --device=/dev/net/tun --dns 209.222.18.222 --dns 209.222.18.218 \
        -e 'REGION=Germany' -v '/yourhostpath/auth.conf:/pia/auth.conf:ro' \
        qmcgaw/private-internet-access
        ```

        Wait about 5 seconds for it to connect to the PIA server.
    1. Check your host IP address with:

        ```bash
        curl -s ifconfig.co
        ```

    1. Run the **curl** Docker container using your *piaTEST* container with:

        ```bash
        docker run --rm --net=container:piaTEST tutum/curl curl -s ifconfig.co
        ```

        If the displayed IP address appears and is different that your host IP address, your PIA OpenVPN client works !    

1. Run the container as a daemon in the background with (and change the `/yourhostpath/auth.conf`):

   ```bash
    docker run -d --restart=always --name=pia --cap-add=NET_ADMIN \
    --device=/dev/net/tun --dns 209.222.18.222 --dns 209.222.18.218 \
    -e 'REGION=Germany' -v '/yourhostpath/auth.conf:/pia/auth.conf' \
    qmcgaw/private-internet-access
    ```
        
### Using Docker Compose

1. Download [**docker-compose.yml**](https://github.com/qdm12/private-internet-access-docker/blob/master/docker-compose.yml)
1. Edit it and change `yourpath`
1. Run the container as a daemon in the background with:

   ```bash
    docker-compose up -d
    ```

    Wait about 5 seconds for it to connect to the PIA server.
1. Check your host IP address with:

    ```bash
    curl -s ifconfig.co
    ```

1. Run the **curl** Docker container using your *pia* container with:

    ```bash
    docker run --rm --net=container:pia tutum/curl curl -s ifconfig.co
    ```

    If the displayed IP address appears and is different that your host IP address, your PIA OpenVPN client works !    


## Connect other containers to it

Connect other Docker containers to the VPN connection by adding `--net=container:pia` when launching them.

## Container launch parameters

- You can change the `--name=` parameter to anything you like
- You can change the `REGION=` parameter to one of the [regions supported by private internet access](https://www.privateinternetaccess.com/pages/network/)
- You must adapt the `/yourhostpath/auth.conf` path to your host path where you created `auth.conf`

## Access ports of containers connected to the VPN container

You have to use another container acting as a Reverse Proxy such as Nginx. 

**Example**:

1. I have a *Deluge* container connected to the PIA container with `--net=container:pia` and its WebUI runs on port 8112.
2. I create the following Nginx configuration file *nginx.conf*:

    ```
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
            listen 80;
            location / {
                proxy_pass http://deluge:8112/;
                proxy_set_header X-Deluge-Base "/";
            }
        }
        include /etc/nginx/conf.d/*.conf;
    }
    ```

3. I run the Alpine [Nginx container](https://hub.docker.com/_/nginx/) with:

    ```bash
    sudo docker -d --restart=always --name=proxypia -p 8000:80 --link pia:deluge \
    -v /mypathto/nginx.conf:/etc/nginx/nginx.conf:ro nginx:alpine
    ```
    
4. Now I can access the WebUI of Deluge at `localhost:8000`
5. You can add more `--link pia:xxx` for more containers and you have to modify *nginx.conf*
