package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/MohamedBassem/servgo/http"
)

var rootDir string

func handleGetRequest(req servgo.Request) servgo.Response {
	res := servgo.NewResponse()
	path := req.Path()[1:]
	if path == "" {
		path = "index.html"
	}
	fullPath := rootDir + path
	absPath, _ := filepath.Abs(fullPath)
	if !strings.HasPrefix(absPath, rootDir) {
		return servgo.NewErrorResponse(403, "403 Forbidden")
	}

	f, err := ioutil.ReadFile(rootDir + path)
	if err != nil {
		return servgo.NewErrorResponse(404, "File not found")
	}

	res.SetBody(string(f))
	res.SetStatusCode(200)

	return res
}

func isValidRootDir(dir string) bool {
	fi, err := os.Stat(dir)
	if err != nil {
		return false
	}
	if os.IsNotExist(err) {
		return false
	}
	return fi.IsDir()
}

func main() {
	var addr = flag.String("addr", ":8080", "The server's port")
	var maxQueuedConnecions = flag.Int("max-queued", 1024, "[Optional] The maximum number of connections that can be queued in the server")
	var numberOfWorkers = flag.Int("num-workers", 2, "[Optional] Number of workers serving the requests")
	var _rootDir = flag.String("root-dir", "nil", "[Required] The root dir for serving the files")
	flag.Parse()
	rootDir = *_rootDir

	if rootDir == "nil" || !isValidRootDir(rootDir) {
		fmt.Println("A --root-dir must be given and valid.")
		os.Exit(1)
	}
	rootDir, _ = filepath.Abs(rootDir)

	if string(rootDir)[len(rootDir)-1:] != "/" {
		rootDir += "/"
	}
	server := servgo.NewServer(*addr, *numberOfWorkers, *maxQueuedConnecions)
	server.SetGetHandler(handleGetRequest)
	err := server.Run()

	if err != nil {
		fmt.Printf("Something went wrong..\n")
	}
}
