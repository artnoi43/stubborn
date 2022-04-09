# 💩 stubborn DNS resolver
I spent my spare time writing this 💩 shitty 💩 DoT/DoH (DNS-over-TLS and DNS-over-HTTPS) resolver just for fun. I have been using [stubby](https://dnsprivacy.org/dns_privacy_daemon_-_stubby/) as my stub DoT resolver for some time, so I just name this program after stubby.

## Features

- DoT/DoH outbound, DNS-over-UDP/53 ("Do53") client replies (default)

- In-memory caching with [patrickmn/go-cache](https://github.com/patrickmn/go-cache)

- Basic yaml configuration (will soon add default config location in `/etc/stubborn/config.yaml`)

- Local network resolver reading host file from `/etc/stubborn/table.json`