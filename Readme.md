# uping

Unidirectional Latency Find Tool  

Determine latency with unidirectional measurements.

This not like run only client side, it is requires server side for measure.

## Server

Avaible server configurations:

- --listen string  
Listen addr for this service (default ":50123")

- --timeserver string  
        Time server which is used for updating time (default "time.cloudflare.com")
- --tsi uint  
        Time sync interval for time client (default 10)

```bash
$ uping server
Server Mode
Time synced, the offset is "365.299µs"
```

## Client

Avaible client configurations:

- -c number  
        Ping count
- -i number  
        Ping interval (default 1)
- -p number  
        Remote server custom port number
- -sp number
        Source port number
- --timeserver string  
        Time server which is used for updating time (default "time.cloudflare.com")
- --tsi number
        Time sync interval for time client (default 10)

```bash
uping client 203.0.113.123
Client Mode
Time synced, the offset is "-307.962038ms"

Pinging from 172.17.0.2:45578 to 203.0.113.123:50123
6 bytes client > 36ms > server(203.0.113.123) > 50ms > client total 86 seq 1
6 bytes client > 36ms > server(203.0.113.123) > 50ms > client total 86 seq 2
6 bytes client > 37ms > server(203.0.113.123) > 49ms > client total 86 seq 3
6 bytes client > 38ms > server(203.0.113.123) > 50ms > client total 88 seq 4
6 bytes client > 38ms > server(203.0.113.123) > 50ms > client total 88 seq 5

5 packets transmitted, 5 received, 0% packet Loss 
Average outbound packet delay 37
Average inbound packet delay 50
```

## Install From Source

For compiling, you have to install go on your system. If you want to use docker, you can also use my "ahmetozer/golang" container for compile and run.

```bash
# Get source
go get github.com/ahmetozer/uping

# To build
go build github.com/ahmetozer/uping

# Run app
uping
```
