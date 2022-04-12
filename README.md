# 💩 stubborn DNS resolver
stubborn is a shitty (💩) DoT/DoH (DNS-over-TLS and DNS-over-HTTPS) resolver that I write just for fun. I have been using [stubby](https://dnsprivacy.org/dns_privacy_daemon_-_stubby/) as my stub DoT + DNSSEC resolver for years, so I borrowed part of stubby name for this project.

## Features

- DoT/DoH outbound, DNS-over-UDP/53 ("Do53") client replies (default)

- In-memory caching with [patrickmn/go-cache](https://github.com/patrickmn/go-cache)

- Basic yaml configuration (will soon add default config location in `/etc/stubborn/config.yaml`)

- Local network resolver reading host file from `/etc/stubborn/table.json` and `/etc/hosts` (UNIX only)

> **IPv6 is only supported in DoH**. There's plan to support IPv6 in DoT too, but not yet implemented. For local network resolver, only IPv4 is supported. But who the fuck uses IPv6 in their house anyway?