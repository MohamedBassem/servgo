# Servgo
A simple HTTP server for static websites written from scratch with thread pools.

### Why?
I decided to write this very simple HTTP server from scratch without using any external packages or the core http library just to learn Go.

### Installation

#### The library
The server doesn't need any external packages, just `go get` this repo and you are ready to go.
```bash
go get github.com/MohamedBassem/servgo
```

#### A Simple HTTP server executable
```bash
go get github.com/MohamedBassem/servgo/servgo
```


### Usage

#### The Library
You can check the executable's source code to know how to use it: 
(https://github.com/MohamedBassem/servgo/blob/master/servgo/main.go)[https://github.com/MohamedBassem/servgo/blob/master/servgo/main.go]

#### The executable
```bash
$ servgo -h
Usage of ./servgo
-addr=":8080": The port to which the server will listen
-max-queued=1024: [Optional] The maximum number of connections that can be queued in the server
-num-workers=2: [Optional] Number of workers serving the requests
-root-dir="nil": [Required] The root dir for serving the files

$ servgo --root-dir example # Will start a server listening to port 8080 and serving files from the example directory
```

### Notes On the Executable
- The server currently supports only `GET` and `HEAD` requests.
- The server spawns `--num-workers` workers to serve the incoming requests.
- If a file is not specified in the request path `index.html` is served.
- Server logs are printed to stdout, if you want to log them to a file you can pipe them.
