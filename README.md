<p>
  <img src="https://github.com/dnsupdater/dnsupdater/raw/main/doc/dns_updater.png">
</p>

# DNS Updater

DNS Updater provides "present" and "cleanup" functions that [**LEGO**](https://go-acme.github.io/lego) expects from the external program provider.

Please refer to our docs at: [https://dnsupdater.github.io](https://dnsupdater.github.io)

LEGO: Let’s Encrypt client and ACME library written in Go.
Please refer to LEGO docs at [https://go-acme.github.io/lego](https://go-acme.github.io/lego)

### DNS Updater usage:

```
dnsu help

# for cPanel provider
 dnsu cpnael --url <cPanel URL> --user <cPanel User> --token <cPanel Token> [--logoutput <log file name>] info <domain name for A record>
 dnsu cpnael --url <cPanel URL> --user <cPanel User> --token <cPanel Token> [--logoutput <log file name>] present <domain name> for TXT record> <auth-key>'
 dnsu cpnael --url <cPanel URL> --user <cPanel User> --token <cPanel Token> [--logoutput <log file name>] cleanup <domain name> for TXT record> <auth-key>'
	
```
### Supported environment eariables:
- DNSU_LOG-OUTPUT for log file
- DNSU_CPANEL-URL for cPanel URL
- DNSU_CPANEL-USER for cPanel User
- DNSU_CPANEL-TOKEN for cPanel Token


### Example:
```
# for verify cPanel access
dnsu cpnael --url "https://cpanel-hostname:2083" --user cpaneluser --token "RMYKKBIT5TQ1ITFU58VZBQB5TDEYQZN4" info '_acme-challenge.my.example.org.'
```

### Example update-dns.sh for LEGO:
```
#!/bin/sh
export DNSU_CPANEL-URL="https://cpanel-hostname:2083"
export DNSU_CPANEL-USER=cpaneluser
export DNSU_CPANEL-TOKEN="RMYKKBIT5TQ1ITFU58VZBQB5TDEYQZN4"
dnsu cpanel "$@"
```


### Create cPanel Manage API Token:
# ![cPanel](./doc/cpanel1.png)


### TODO:
- [**LEGO**](https://go-acme.github.io/lego/dns/httpreq/) HTTP request provider.
