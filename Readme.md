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

- -i float  
        Ping interval (default 1)
- --timeserver string  
        Time server which is used for updating time (default "time.cloudflare.com")
- --tsi uint  
        Time sync interval for time client (default 10)

```bash
uping client 203.0.113.123
Client Mode
Time synced, the offset is "-307.962038ms"
6 bytes client > 36ms > server(203.0.113.123) > 50ms > client total 86
6 bytes client > 36ms > server(203.0.113.123) > 50ms > client total 86
6 bytes client > 37ms > server(203.0.113.123) > 49ms > client total 86
6 bytes client > 38ms > server(203.0.113.123) > 50ms > client total 88
6 bytes client > 38ms > server(203.0.113.123) > 50ms > client total 88
```

## Install From Source

For compiling, you have to install go on your system. If you use docker, you can also use my "ahmetozer/golang" container for compile.

You can get source with git or download zip from github.

```bash
git clone --depth 1 --single-branch  git@github.com:ahmetozer/uping.git

# Get required libraries to build.
go get -v .

# To build
go build

# Move to "/usr/bin/uping"
mv uping  /usr/bin/uping
```
