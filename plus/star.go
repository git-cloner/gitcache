package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

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

func GetRepoStar(w http.ResponseWriter, r *http.Request) {
	github_token := os.Getenv("GITHUB_TOKEN")
	url := "https://api.github.com/repos/" + strings.Replace(r.URL.RequestURI(), "/gitcache/star/", "", -1) + "?access_token=" + github_token
	log.Printf("get star url : %v \n", url)
	contents := httpGet(url)
	const reg = `"stargazers_count":\s*(\d+),`
	compile := regexp.MustCompile(reg)
	submatch := compile.FindAllSubmatch([]byte(contents), -1)
	var stargazers_count = 0
	for _, m := range submatch {
		stargazers_count, _ = strconv.Atoi(string(m[1]))
	}
	w.WriteHeader(200)
	w.Write([]byte(strconv.Itoa(stargazers_count)))
	log.Printf("get star count : %v \n", stargazers_count)
}
