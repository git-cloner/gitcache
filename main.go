package main

import (
	"flag"
	"log"
	"net/http"
)

var g_Basedir string
var port string
var global_ssh string

func main() {
	//log params
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.SetPrefix("LOG: ")

	//flag params
	flag.StringVar(&g_Basedir, "b", "/var/gitcache", "default path: /var/gitcache")
	flag.StringVar(&port, "p", "5000", "default port:5000")
	flag.StringVar(&global_ssh, "ssh", "0", "default ssh:0")
	//if set -ssh 1 ,please
	//eval $(ssh-agent -s)
	//ssh-add ~/.ssh/id_rsa
	flag.Parse()

	log.Printf("cache basedir:%v , port:%v, ssh:%v", g_Basedir, port, global_ssh)
	//port == 5000 gitcacher , posrt != 5000 downloader
	//connect to db
	InitDb()
	//cron
	Cron()
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
