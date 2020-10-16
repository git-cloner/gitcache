package main

import (
	"log"
	"net/http"

	"github.com/elazarl/goproxy"
)

func main() {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true
	log.Println("git accelerater listen at 9999")
	log.Fatal(http.ListenAndServe(":9999", proxy))
}
