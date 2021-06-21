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

## Misc Options
This tool supports misc options which are configured under the `misc` array inside of the config file. The only event supported for this at the moment is when a server is restarted from the tool. However, other events may be added in the future. An example may be found below.

```JSON
{
        "misc": [
                {
                        "type": "misctype",
                        "data": {
                                "option1": "val1",
                                "option2": "val2"
                        }
                }
        ]
}
```

### Discord Web Hooks
As of right now, the only misc option `type` is `webhook` which indicates a Discord web hook. Please look at the following data items:

* `id` => The Discord web hook ID (**REQUIRED**).
* `token` => The Discord web hook token (**REQUIRED**).
* `contents` => The contents of the web hook.
* `username` => The username the web hook send as.
* `avatarurl` => The avatar URL used with the web hook.

**Note** - Retrieve the web hook's token and ID from the full web hook URL as the following:

```
https://discord.com/api/webhooks/<id>/<token>
```

#### Variable Replacements For Contents
The following strings are replaced inside of the `contents` string before the web hook submission.

* `{IP}` => The server's IP.
* `{Port}` => The server's port.
* `{FAILS}` => The server's current fail count.
* `{RESTARTS}` => The amount of times the server has been restarted since down.
* `{MAXFAILS}` => The server's configured max fails.
* `{MAXRESTARTS}` => The server's configured max restarts.
* `{UID}` => The server's UID from the config file/Pterodactyl API.
* `{SCANTIME}` => The server's configured scan time.
* `{RESTARTINT}` => The server's configured restart interval.

#### Defaults
Here are the Discord web hook's default values.

* `contents` => \*\*SERVER DOWN\*\*\\n- IP => {IP}\\n- Port => {PORT}\\n- Fail Count => {FAILS}/{MAXFAILS}
* `username` => Pterowatch
* `avatarurl` => *empty* (default)

## Configuration Example
Here's an configuration example in JSON:

```JSON
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

You may find other config examples in the [tests/](https://github.com/gamemann/Pterodactyl-Game-Server-Watch/tests) directory.

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