package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/robfig/cron"
)

type LocalMirrorsInfo struct {
	Count    int64  `json:"count"`
	Progress string `json:"progress"`
	Size     int64  `json:"size"`
	Nodes    string `json:"nodes"`
}

var _PATH_DEPTH = 2
var _IS_SYNC = false
var _REPO_COUNT int64 = 0

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

func countCacheRepository(repository string) {
	_REPO_COUNT++
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
	_REPO_COUNT = 0
	walkDir(g_Basedir, 0, countCacheRepository)
	info := LocalMirrorsInfo{}
	info.Count = _REPO_COUNT
	info.Nodes = ""
	info.Progress = ""
	info.Size = 0
	data, _ := json.Marshal(info)
	return string(data)
}

func httpPost(url string, contentType string, body string) string {
	resp, err := http.Post(url, contentType, strings.NewReader(body))
	if err != nil {
		return err.Error()
	}
	defer resp.Body.Close()
	rbody, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		return err1.Error()
	}
	return string(rbody)
}

func BroadCastGitCloneCommandToChain(repository string) {
	log.Println("broadcast git clone command to chain : " + repository)
	var body = "{\"privatekey\":\"f45b1d6e433195a0e70a09ffaf59d4c71bc9c49f84cfe63fd455b3c34a8fcd2d246ea5c7d47cf6027e4ec99b27dade8e23bb811a07b90228c3f27f744c4d1322\"," +
		"\"publickey\":\"" + "246EA5C7D47CF6027E4EC99B27DADE8E23BB811A07B90228C3F27F744C4D1322\"," +
		"\"msg\":\"" + "git clone " + repository + "\"}"
	go httpPost("http://172.16.62.48:4000/broadcast/msg", "application/json", body)
}

func Cron() {
	c := cron.New()
	c.AddFunc("0 0 0,2,4,6,20,22 * * *", func() {
		//c.AddFunc("0 */1 * * * *", func() {
		SyncLocalMirrorFromRemote()
	})
	c.Start()
	log.Println("cron start")
	return
}
