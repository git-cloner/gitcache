package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
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
var _REPO_ALL_COUNT int64 = 0

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

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}

func countCacheRepositoryByIP(url string) int64 {
	var group_repos_info string
	var ct int64 = 0
	group_repos_info = httpGet(url)
	if group_repos_info != "" {
		var localMirrorsInfo LocalMirrorsInfo
		json.Unmarshal([]byte(group_repos_info), &localMirrorsInfo)
		ct = localMirrorsInfo.Count
		log.Printf("%+v:%+v%+v\n", url, ct, group_repos_info)

	}
	return ct
}

func countAllCacheRepository() int64 {
	var ct int64
	ct = countCacheRepositoryByIP("http://192.168.10.54:5000/gitcache/system/info")
	log.Printf("54:%+v\n", ct)
	ct = ct + countCacheRepositoryByIP("http://192.168.10.55:5000/gitcache/system/info")
	log.Printf("55:%+v\n", ct)
	ct = ct + countCacheRepositoryByIP("http://192.168.10.56:5000/gitcache/system/info")
	log.Printf("56:%+v\n", ct)
	ct = ct + countCacheRepositoryByIP("http://192.168.10.57:5000/gitcache/system/info")
	log.Printf("57:%+v\n", ct)
	return ct
}

func SyncCountCacheRepository() {
	_REPO_COUNT = 0
	walkDir(g_Basedir, 0, countCacheRepository)
	_REPO_ALL_COUNT = countAllCacheRepository()
	log.Printf("sync count cache repository : %v\n", _REPO_COUNT)
}

func httpGet(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return string(body)
}

func GetLocalMirrorsInfo() string {
	if _REPO_COUNT == 0 {
		walkDir(g_Basedir, 0, countCacheRepository)
	}
	info := LocalMirrorsInfo{}
	info.Count = _REPO_COUNT
	var ip net.IP = GetOutboundIP()
	var sip = "node:0"
	str := strings.Split(ip.String(), ".")
	if len(str) == 4 {
		sip = "node:" + str[3]
	}
	info.Nodes = sip
	info.Progress = ""
	if _REPO_ALL_COUNT > 0 {
		info.Size = _REPO_ALL_COUNT
	} else {
		info.Size = _REPO_COUNT
	}
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
	var msgtx MsgTx
	msgtx.PrivateKey = "f45b1d6e433195a0e70a09ffaf59d4c71bc9c49f84cfe63fd455b3c34a8fcd2d246ea5c7d47cf6027e4ec99b27dade8e23bb811a07b90228c3f27f744c4d1322"
	msgtx.PublicKey = "246EA5C7D47CF6027E4EC99B27DADE8E23BB811A07B90228C3F27F744C4D1322"
	msgtx.Msg = "git clone " + repository
	go BroadCastMsg(msgtx)
}

func Cron() {
	c := cron.New()
	c.AddFunc("0 0 0,4,12,20 * * *", func() {
		go SyncLocalMirrorFromRemote()
	})
	c.AddFunc("0 */1 * * * *", func() {
		go SyncCountCacheRepository()
	})
	c.Start()
	log.Println("cron start")
	return
}
