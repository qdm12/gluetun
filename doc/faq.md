# FAQ

## Table of content

- [Openvpn disconnects because of a ping timeout](#Openvpn-disconnects-because-of-a-ping-timeout)
- [Private Internet Access: Why do I see openvpn warnings at start](#Private-Internet-Access:-Why-do-I-see-openvpn-warnings-at-start)
- [What files does it download after tunneling](#What-files-does-it-download-after-tunneling)
- [How to build Docker images of older or alternate versions](#How-to-build-Docker-images-of-older-or-alternate-versions)
- [Mullvad does not work with IPv6](#Mullvad-does-not-work-with-IPv6)
- [What's all this Go code](#What-is-all-this-Go-code)
- [How to test DNS over TLS](#How-to-test-DNS-over-TLS)

## Openvpn disconnects because of a ping timeout

It happens especially on some PIA servers where they change their configuration or the server goes offline.

You will obtain an error similar to:

```s
openvpn: Wed Mar 18 22:13:00 2020 [3a51ae90324bcb0719cb399b650c64d4] Inactivity timeout (--ping-restart), restarting,
openvpn: Wed Mar 18 22:13:00 2020 SIGUSR1[soft,ping-restart] received, process restarting,
...
openvpn: Wed Mar 18 22:13:17 2020 Preserving previous TUN/TAP instance: tun0,
openvpn: Wed Mar 18 22:13:17 2020 NOTE: Pulled options changed on restart, will need to close and reopen TUN/TAP device.,
openvpn: Wed Mar 18 22:13:17 2020 ERROR: Linux route delete command failed: external program exited with error status: 2,
openvpn: Wed Mar 18 22:13:17 2020 ERROR: Linux route delete command failed: external program exited with error status: 2,
openvpn: Wed Mar 18 22:13:17 2020 ERROR: Linux route delete command failed: external program exited with error status: 2,
openvpn: Wed Mar 18 22:13:17 2020 ERROR: Linux route delete command failed: external program exited with error status: 2,
openvpn: Wed Mar 18 22:13:17 2020 /sbin/ip addr del dev tun0 local 10.6.11.6 peer 10.6.11.5,
openvpn: Wed Mar 18 22:13:17 2020 Linux ip addr del failed: external program exited with error status: 2,
openvpn: Wed Mar 18 22:13:18 2020 ERROR: Cannot ioctl TUNSETIFF tun: Operation not permitted (errno=1),
openvpn: Wed Mar 18 22:13:18 2020 Exiting due to fatal error,
exit status 1
```

To fix it, you would have to run openvpn with root, by setting the environment variable `OPENVPN_ROOT=yes`.

## Private Internet Access: Why do I see openvpn warnings at start

You might see some warnings similar to:

```s
openvpn: Sat Feb 22 15:55:02 2020 WARNING: this configuration may cache passwords in memory -- use the auth-nocache option to prevent this
openvpn: Sat Feb 22 15:55:02 2020 WARNING: 'link-mtu' is used inconsistently, local='link-mtu 1569', remote='link-mtu 1542'
openvpn: Sat Feb 22 15:55:02 2020 WARNING: 'cipher' is used inconsistently, local='cipher AES-256-CBC', remote='cipher BF-CBC'
openvpn: Sat Feb 22 15:55:02 2020 WARNING: 'auth' is used inconsistently, local='auth SHA256', remote='auth SHA1'
openvpn: Sat Feb 22 15:55:02 2020 WARNING: 'keysize' is used inconsistently, local='keysize 256', remote='keysize 128'
openvpn: Sat Feb 22 15:55:02 2020 WARNING: 'comp-lzo' is present in remote config but missing in local config, remote='comp-lzo'
openvpn: Sat Feb 22 15:55:02 2020 [a121ce520d670b71bfd3aa475485539b] Peer Connection Initiated with [AF_INET]xx.xx.xx.xx:1197
```

It is mainly because the option [disable-occ](https://openvpn.net/community-resources/reference-manual-for-openvpn-2-4/) was removed for transparency with you.

Private Internet Access explains [here why](https://www.privateinternetaccess.com/helpdesk/kb/articles/why-do-i-get-cipher-auth-warnings-when-i-connect) the warnings show up.

## What files does it download after tunneling

At start, after tunneling, the Go entrypoint only downloads, depending on your settings:

- If `DOT=on`: [DNS over TLS named root](https://github.com/qdm12/files/blob/master/named.root.updated) for Unbound
- If `DOT=on`: [DNS over TLS root key](https://github.com/qdm12/files/blob/master/root.key.updated) for Unbound
- If `BLOCK_MALICIOUS=on`: [Malicious hostnames and IP addresses block lists](https://github.com/qdm12/files) for Unbound
- If `BLOCK_SURVEILLANCE=on`: [Surveillance hostnames and IP addresses block lists](https://github.com/qdm12/files) for Unbound
- If `BLOCK_ADS=on`: [Ads hostnames and IP addresses block lists](https://github.com/qdm12/files) for Unbound

## How to build Docker images of older or alternate versions

First, install [Git](https://git-scm.com/).

The following will build the Docker image locally and replace the previous one you built or pulled.

- Build the latest image

    ```sh
    docker build -t qmcgaw/private-internet-access https://github.com/qdm12/private-internet-access-docker.git
    ```

- Or, find a [commit](https://github.com/qdm12/private-internet-access-docker/commits/master) you want to build for, in example `095623925a9cc0e5cf89d5b9b510714792267d9b`, then:

    ```sh
    docker build -t qmcgaw/private-internet-access https://github.com/qdm12/private-internet-access-docker.git#095623925a9cc0e5cf89d5b9b510714792267d9b
    ```

- Or, find a [branch](https://github.com/qdm12/private-internet-access-docker/branches) you want to build for, in example `mullvad`, then:

    ```sh
    docker build -t qmcgaw/private-internet-access https://github.com/qdm12/private-internet-access-docker.git#mullvad
    ```

## Mullvad does not work with IPv6

By default, the Mullvad server tunnels both ipv4 and ipv6, hence openvpn will try to create an
ipv6 route. To allow the container to create such route, you have to specify `net.ipv6.conf.all.disable_ipv6=0`
at runtime, using either:

- For a Docker run command, the flag: `--sysctl net.ipv6.conf.all.disable_ipv6=0`
- In a docker-compose file:

    ```yml
        sysctls:
          - net.ipv6.conf.all.disable_ipv6=0
    ```

## What is all this Go code

The Go code is a big rewrite of the previous shell entrypoint, it allows for:

- better testing
- better maintainability
- ease of implementing new features
- faster boot
- asynchronous/parallel operations

It is mostly made of the [internal directory](../internal) and the entry Go file [cmd/main.go](../cmd/main.go).

## How to test DNS over TLS

- You can test DNSSEC using [internet.nl/connection](https://www.internet.nl/connection/)
- Check DNS leak tests with [https://www.dnsleaktest.com](https://www.dnsleaktest.com)
- Some other DNS leaks tests might not work because of [this](https://github.com/qdm12/cloudflare-dns-server#verify-dns-connection) (*TLDR*: Unbound DNS server is a local caching intermediary)
