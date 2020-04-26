package main

import (
	"io/ioutil"
	"log"
	"strings"

	"github.com/robfig/cron"
)

var _PATH_DEPTH = 2
var _IS_SYNC = false

func fetchMirrorFromRemoteUnshallow(repository string) {
	remote := "https:/" + strings.Replace(repository, g_Basedir, "", -1)
	local := repository
	log.Printf("git remote update: %s begin\n", local)
	err := fetchMirrorFromRemote(remote, local, "update")
	if err == "" {
		err = "ok"
	}
	log.Printf("git remote update: %s %s\n", local, err)
}

func walkDir(dirpath string, depth int, f func(string)) {
	if depth > _PATH_DEPTH {
		return
	}
	files, err := ioutil.ReadDir(dirpath)
	if err != nil {
		return
	}
	for _, file := range files {
		if file.IsDir() {
			walkDir(dirpath+"/"+file.Name(), depth+1, f)
			headExist, _ := PathExists(dirpath + "/" + file.Name() + "/HEAD")
			if headExist && (!strings.HasSuffix(file.Name(), "logs")) {
				f(dirpath + "/" + file.Name())
			}
			continue
		}
	}
}

func SyncLocalMirrorFromRemote() {
	if _IS_SYNC {
		log.Println("syncing local mirror from remote,sync ignore")
		return
	}
	log.Println("sync local mirror from remote begin")
	_IS_SYNC = true
	walkDir(g_Basedir, 0, fetchMirrorFromRemoteUnshallow)
	log.Println("sync local mirror from remote end")
	_IS_SYNC = false
}

func GetLocalMirrorsInfo() string {
	return ""
}

func GetMirrorProgress(repoName string) string {
	return ""
}

func Cron() {
	c := cron.New()
	c.AddFunc("0 0 */2 * * *", func() {
		//c.AddFunc("0 */1 * * * *", func() {
		SyncLocalMirrorFromRemote()
	})
	c.Start()
	log.Println("cron start")
	return
}
