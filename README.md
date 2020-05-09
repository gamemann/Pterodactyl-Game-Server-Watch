# Pterodactyl Game Server Watch

## Description
A tool programmed in Go to automatically restart 'hung' game servers/containers via a Pterodactyl API (version 0.7). This only supports game servers that respond to the [A2S_INFO](https://developer.valvesoftware.com/wiki/Server_queries#A2S_INFO) query (a Valve Master Server query).

## Config File
The config file's default path is `/etc/pterowatch/pterowatch.conf`. This should be a JSON array including the API URL, token, and an array of servers to check against. The main options are the following:

* `apiURL` => The Pterodactyl API URL.
* `token` => The bearer token to use when sending HTTP POST requests to the Pterodactyl API.
* `servers` => An array of servers to check against (read below).

The `servers` array should contain the following members:

* `enable` => If true, this server will be scanned.
* `ip` => The IP to send A2S_INFO requests to.
* `port` => The port to send A2S_INFO requests to.
* `uid` => The server's Pterodactyl UID.
* `scantime` => How often to scan a game server/container in seconds.
* `maxfails` => The maximum amount of A2S_INFO response failures before attempting to restart the game server/container.
* `maxrestarts` => The maximum amount of times we attempt to restart the server until A2S_INFO responses start coming back successfully.
* `restartint` => When a game server/container is restarted, the program won't start scanning the server until *x* seconds later.

## Configuration Example
Here's an configuration example in JSON:

```
{
        "apiURL": "https://panel.mydomain.com",
        "token": "12345",

        "servers": [
                {
                        "enable": true,
                        "ip": "172.20.0.10",
                        "port": 27015,
                        "uid": "testingUID",
                        "scantime": 5,
                        "maxfails": 5,
                        "maxrestarts": 1,
                        "restartint": 120
                },
                {
                        "enable": true,
                        "ip": "172.20.0.11",
                        "port": 27015,
                        "uid": "testingUID2",
                        "scantime": 5,
                        "maxfails": 10,
                        "maxrestarts": 2,
                        "restartint": 120
                }
        ]
}
```

## Building
You may use `git` and `go build` to build this project and produce a binary. Example:

```
git clone https://github.com/gamemann/Pterodactyl-Game-Server-Watch.git
cd Pterodactyl-Game-Server-Watch/src
go build
```

## Credits
* [Christian Deacon](https://www.linkedin.com/in/christian-deacon-902042186/) - Creator.