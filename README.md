#Juhua dns proxy

This is a dns proxy to fuck `G.F.W` by golang.

##Default dns server to lookup domain

- 8.8.8.8 Google
- 8.8.4.4 Google
- 156.154.70.1 Dnsadvantage
- 156.154.71.1 Dnsadvantage
- 208.67.222.222 OpenDNS
- 208.67.220.220 OpenDNS
- 198.153.192.1 Norton
- 198.153.194.1 Norton

Default load all of above severs.

##How to use it?

Install it via golang tools and start proxy

Change your /etc/resolve.conf

> mv /etc/resolve.conf /etc/resolve.conf.bak
> echo "nameserver 127.0.0.1" > /etc/resolve.conf

