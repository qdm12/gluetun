---
name: Support a VPN provider
about: Suggest a VPN provider to be supported
title: 'VPN provider support: NAME OF THE PROVIDER'
labels: ":bulb: New provider"

---

Important notes:

- There is no need to support both OpenVPN and Wireguard for a provider, but it's better to support both if possible
- We do **not** implement authentication to access servers information behind a login. This is way too time consuming unfortunately
- If it's not possible to support a provider natively, you can still use the [the custom provider](https://github.com/qdm12/gluetun-wiki/blob/main/setup/providers/custom.md)

## For Wireguard

Wireguard can be natively supported ONLY if:

- the `PrivateKey` field value is the same across all servers for one user account
- the `Address` field value is:
  - can be found in a structured (JSON etc.) list of servers publicly available; OR
  - the same across all servers for one user account
- the `PublicKey` field value is:
  - can be found in a structured (JSON etc.) list of servers publicly available; OR
  - the same across all servers for one user account
- the `Endpoint` field value:
  - can be found in a structured (JSON etc.) list of servers publicly available
  - can be determined using a pattern, for example using country codes in hostnames

If any of these conditions are not met, Wireguard cannot be natively supported or there is no advantage compared to using a custom Wireguard configuration file.

If **all** of these conditions are met, please provide an answer for each of them.

## For OpenVPN

OpenVPN can be natively supported ONLY if one of the following can be provided, by preference in this order:

- Publicly accessible URL to a structured (JSON etc.) list of servers **and attach** an example Openvpn configuration file for both TCP and UDP; OR
- Publicly accessible URL to a zip file containing the Openvpn configuration files; OR
- Publicly accessible URL to the list of servers **and attach** an example Openvpn configuration file for both TCP and UDP
