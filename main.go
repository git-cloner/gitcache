package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var g_Basedir string
var port string
var global_ssh string
var global_blacklist []string

func getCurrentAbPathByExecutable() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	res, _ := filepath.EvalSymlinks(filepath.Dir(exePath))
	return res
}

func ReadBlacklist() {
	filepath := getCurrentAbPathByExecutable() + "/blacklist.txt"
	_, err := os.Stat(filepath)
	if err == nil {
		s, err := ioutil.ReadFile(filepath)
		if err != nil {
			log.Fatal(err)
		}
		global_blacklist = strings.Split(string(s), "\r\n")
	} else {
		global_blacklist = []string{}
	}
	log.Printf("blacklist:%v", global_blacklist)
}

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
	//read blacklist
	ReadBlacklist()
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
