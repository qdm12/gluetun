# Private Internet Access Docker (OpenVPN, Alpine)

[![PIA Docker OpenVPN](https://github.com/qdm12/private-internet-access-docker/raw/master/readme/title.png)](https://hub.docker.com/r/qmcgaw/private-internet-access/)

VPN client container to private internet access servers based on [Alpine Linux](https://alpinelinux.org/) and [OpenVPN](https://openvpn.net/)

It requires:
- A Private Internet Access **username** and **password** - [signup up](https://www.privateinternetaccess.com/pages/buy-vpn/)
- [Docker](https://docs.docker.com/install/) installed on the host

The PIA configuration files are downloaded from [the PIA website](https://www.privateinternetaccess.com/openvpn/openvpn.zip) when the Docker image gets built.

## Installation & Testing

1. Run the *tun.sh* script on your host machine to ensure you have the **tun** device setup

    ```bash
    sudo chmod +x tun.sh
    ./tun.sh
    ```

2. Obtaining the Docker image
    - Option 1 of 2: Automated download from the Docker Hub Registry, simply go to step 3.
    - Option 2 of 2: Build the image
        1. Download the repository files or `git clone` them
        2. With a terminal, go in the directory where the *Dockerfile* is located
        3. Build the image with:

            ```bash
            sudo docker build -t qmcgaw/private-internet-access ./
            ```

3. Create a file *auth.conf* in `/yourhostpath` (for example), with:
    - On the first line: your PIA username (i.e. `js89ds7`)
    - On the second line: your PIA password (i.e. `8fd9s239G`)
4. Test the container by connecting another container to it
    1. Run the container interactively with (and change the `/yourhostpath/auth.conf`):

        ```bash
        sudo docker run --rm --name=piaTEST --cap-add=NET_ADMIN \
        --device=/dev/net/tun --dns 209.222.18.222 --dns 209.222.18.218 \
        -e 'REGION=Romania' -v '/yourhostpath/auth.conf:/pia/auth.conf' \
        qmcgaw/private-internet-access
        ```

        Wait about **5** seconds for it to connect to the PIA server.
    2. Check your host IP address with:

        ```bash
        curl -s ifconfig.co
        ```

    3. Run the **curl** Docker container using your *piaTEST* container with:

        ```bash
        sudo docker run --rm --net=container:piaTEST tutum/curl curl -s ifconfig.co
        ```

        If the displayed IP address appears and is different that your host IP address, your PIA OpenVPN client works !    
5. Run the container as a daemon in the background with (and change the `/yourhostpath/auth.conf`):

   ```bash
    sudo docker run -d --restart=always --name=pia --cap-add=NET_ADMIN \
    --device=/dev/net/tun --dns 209.222.18.222 --dns 209.222.18.218 \
    -e 'REGION=Romania' -v '/yourhostpath/auth.conf:/pia/auth.conf' \
    qmcgaw/private-internet-access
    ```

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
