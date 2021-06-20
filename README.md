# Pterodactyl Game Server Watch

## Description
A tool programmed in Go to automatically restart 'hung' (game) servers via the Pterodactyl API (working since version 1.4.2). This only supports servers that respond to the [A2S_INFO](https://developer.valvesoftware.com/wiki/Server_queries#A2S_INFO) query (a Valve Master Server query). I am currently looking for a better way to detect if a server is hung, though.

## Config File
The config file's path is `/etc/pterowatch/pterowatch.conf`. This should be a JSON array including the API URL, token, and an array of servers to check against. The main options are the following:

* `apiurl` => The Pterodactyl API URL (do not include the `/` at the end).
* `token` => The bearer token to use when sending requests to the Pterodactyl API.
* `addservers` => Whether or not to automatically add servers to the config from the Pterodactyl API.
* `servers` => An array of servers to check against (read below).

The `servers` array should contain the following items:

* `enable` => If true, this server will be scanned.
* `ip` => The IP to send A2S_INFO requests to.
* `port` => The port to send A2S_INFO requests to.
* `uid` => The server's Pterodactyl UID.
* `scantime` => How often to scan a game server in seconds.
* `maxfails` => The maximum amount of A2S_INFO response failures before attempting to restart the game server.
* `maxrestarts` => The maximum amount of times we attempt to restart the server until A2S_INFO responses start coming back successfully.
* `restartint` => When a game server is restarted, the program won't start scanning the server until *x* seconds later.

## Egg Variable Overrides
If you have the `addservers` setting set to true (servers are automatically retrieved via the Pterodactyl API), you may use the following egg variables as overrides to the specific server's config.

* `PTEROWATCH_DISABLE` => If set to above 0, will disable the specific server from the tool.
* `PTEROWATCH_IP` => If not empty, will override the server IP to scan with this value for the specific server.
* `PTEROWATCH_PORT` => If not empty, will override the server port to scan with this value for the specific server.
* `PTEROWATCH_SCANTIME` => If not empty, will override the scan time with this value for the specific server.
* `PTEROWATCH_MAXFAILS` => If not empty, will override the maximum fails with this value for the specific server.
* `PTEROWATCH_MAXRESTARTS` => If not empty, will override the maximum restarts with this value for the specific server.
* `PTEROWATCH_RESTARTINT` => If not empty, will override the restart interval with this value for the specific server.

## Configuration Example
Here's an configuration example in JSON:

```
{
        "apiurl": "https://panel.mydomain.com",
        "token": "12345",
        "addservers": true,

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
You may use `git` and `go build` to build this project and produce a binary. I'd suggest cloning this to `$GOPATH` so there aren't problems with linking modules. For example:

```
cd <Path To One $GOPATH>
git clone https://github.com/gamemann/Pterodactyl-Game-Server-Watch.git
cd Pterodactyl-Game-Server-Watch
go build
```

## Credits
* [Christian Deacon](https://github.com/gamemann) - Creator.