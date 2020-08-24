package main

import (
	"flag"
	"log"
	"net/http"
)

var g_Basedir string
var port string

func main() {
	//log params
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.SetPrefix("LOG: ")
	//flag params
	flag.StringVar(&g_Basedir, "b", "/var/gitcache", "default path: /var/gitcache")
	flag.StringVar(&port, "p", "5000", "default port:5000")
	flag.Parse()
	log.Printf("cache basedir:%v , port:%v", g_Basedir, port)
	//port == 5000 gitcacher , posrt != 5000 downloader
	if port == "5000" {
		//connect to db
		InitDb()
		//cron
		Cron()
	}
	//listen
	http.HandleFunc("/", RequestHandler(g_Basedir))
	address := ":" + port
	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal("ListenAndServer: ", err)
	} else {
		log.Printf("ListenAndServer: %s", address)
	}
}
