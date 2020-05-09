# GFL-Watch

## Description
A tool programmed in Go to automatically restart 'hung' game servers/containers via a Pterodactyl API. This only supports game servers that respond to the A2S_INFO query (a Valve Master Server query).

## Config File
The config file's default path is `/etc/gflwatch/gflwatch.conf`. This should be a JSON array including the API URL, token, and an array of servers to check against. The main options are the following:

* `apiURL` => The Pterodactyl API URL.
* `token` => The bearer token to use when sending HTTP POST requests to the Pterodactyl API.
* `servers` => An array of servers to check against (read below).

The `servers` array should contain the following members:

* `IP` => The IP to send A2S_INFO requests to.
* `port` => The port to send the A2S_INFO requests to.
* `uid` => The server's Pterodactyl UID.

## Configuration Example
Here's an configuration example in JSON:

```
{
        "apiURL": "https://panel.mydomain.com",
        "token": "12345",

        "servers": [
                {
                        "IP": "172.20.0.10",
                        "port": 27015,
                        "uid": "testingUID"
                }
        ]
}
```

## Status
Not finished.

## Credits
* [Christian Deacon](https://www.linkedin.com/in/christian-deacon-902042186/) - Creator.