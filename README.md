# GoServe
A simple HTTP server for static websites written from scratch with thread pools.

### Why?
I decided to write this very simple HTTP server from scratch without using any external packages or the core http library just to learn Go.

### Installation
```bash
go get github.com/MohamedBassem/goserve
```

### Usage
```bash
$ ./goserve -h
Usage of ./goserve:
-addr=":8080": The server's port
-max-queued=1024: [Optional] The maximum number of connections that can be queued in the server
-num-workers=2: [Optional] Number of workers serving the requests
-root-dir="nil": [Required] The root dir for serving the files

$ ./goserve --root-dir example # Will start a server listening to port 8080 and serving files from the example directory
```

### Thread Pools
The server spawns `--num-workers` workers to serve the incoming requests.

### Notes
- The server currently supports only `GET` and `HEAD` requests.
