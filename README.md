# 💩 stubborn DNS resolver
I spent my spare time writing this 💩 shitty 💩 DoH (DNS-over-HTTPS) resolver just for fun. I have been using [stubby](https://dnsprivacy.org/dns_privacy_daemon_-_stubby/) as my stub DoT resolver for some time, so I just name this program after stubby.

## Features

- DoH outbound, DNS-over-UDP/53 ("Do53") client replies (default)

- In-memory caching with [patrickmn/go-cache](https://github.com/patrickmn/go-cache)

- Redis read/write support in case you want to share these resolves with other services

- Basic yaml configuration (will soon add default config location in `$HOME/.config`)