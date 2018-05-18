package main

import (
	"fmt"
	"flag"
	"log"
	"github.com/golang/glog"
	"net/http"
)

func main() {
    serverPort := flag.String("port", "8000", "http port number where server will run")
    flag.Parse()
		blockchain := CreateGennsisChain()
    http.Handle("/", NewHandler(blockchain))
		log.Printf("Starting gochain HTTP Server. Listening at port %q", *serverPort)
		glog.Info("server start....")
    http.ListenAndServe(fmt.Sprintf(":%s", *serverPort), nil)
}
