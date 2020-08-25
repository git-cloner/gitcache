package main

import (
	"flag"
	"log"
	"net/http"
)

var port string

func main() {
	//log params
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.SetPrefix("LOG: ")
	//flag params
	flag.StringVar(&port, "p", "5000", "default port:5000")
	flag.Parse()
	log.Printf("port:%v", port)
	//listen
	http.HandleFunc("/", RequestHandler())
	address := ":" + port
	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal("ListenAndServer: ", err)
	} else {
		log.Printf("ListenAndServer: %s", address)
	}
}
