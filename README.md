# Proxy Redirector

A HTTP/S proxy that redirects connections.
Designed to be used as system proxy or forced for specific programs via software like [Proxifier](https://www.proxifier.com/).

## Help screen (run with `-help`)

```
  -help
        Help screen
  -port int
        Port to listen on (default 8868)
  -settings string
        Settings file with routes (default "settings.txt")
  -verbose
        Verbose proxy output (default true)

Settings file consists of lines defining redirection routes.

A redirection route has the following format:
{ip}:{port} {ip}:{port}

Multiple whitespaces/tabs are permitted as a separator.

This program will redirect the first (source) ip&port to the second (destination) ip&port.
You can use a wildcard '*' to match the ip, port, or both, for the source.

Sample settings.txt:
1.2.3.4:80  127.0.0.1:8080
*:1555      127.0.0.1:6000
*:*         127.0.0.1:80
```