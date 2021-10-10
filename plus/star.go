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

func httpGet(url string, token string) string {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
		return ""
	}
	req.Header.Set("Authorization", "token "+token)
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		log.Println(err)
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return string(body)
}

func GetRepoStar(w http.ResponseWriter, r *http.Request) {
	//sample: https://api.github.com/repos/git-cloner/gitcache
	//test : http://127.0.0.1:5001/gitcache/star/git-cloner/gitcache
	github_token := os.Getenv("GITHUB_TOKEN")
	url := "https://api.github.com/repos/" + strings.Replace(r.URL.RequestURI(), "/gitcache/star/", "", -1)
	log.Printf("get star url : %v \n", url)
	contents := httpGet(url, github_token)
	//get stat
	const reg = `"stargazers_count":\s*(\d+),`
	compile := regexp.MustCompile(reg)
	submatch := compile.FindAllSubmatch([]byte(contents), -1)
	var stargazers_count = 0
	for _, m := range submatch {
		stargazers_count, _ = strconv.Atoi(string(m[1]))
	}
	//get language
	const reg_lang = `\"language\"\:\s*\"(\w+)\",`
	compile_lang := regexp.MustCompile(reg_lang)
	submatch_lang := compile_lang.FindAllSubmatch([]byte(contents), -1)
	var language = ""
	for _, m := range submatch_lang {
		language = string(m[1])
	}
	//get description
	const reg_desc = `\"description\"\:\s*\"([^\"]+)\",`
	compile_desc := regexp.MustCompile(reg_desc)
	submatch_desc := compile_desc.FindAllSubmatch([]byte(contents), -1)
	var description = ""
	for _, m := range submatch_desc {
		description = string(m[1])
	}
	//get updated_at
	const reg_update = `\"updated_at\"\:\s*\"([^\"]+)\",`
	compile_update := regexp.MustCompile(reg_update)
	submatch_update := compile_update.FindAllSubmatch([]byte(contents), -1)
	var updated_at = ""
	for _, m := range submatch_update {
		updated_at = string(m[1])
	}
	w.WriteHeader(200)
	json := `{"stargazers_count":` + strconv.Itoa(stargazers_count) + `,"language":"` + language +
		`","description":"` + description + `","updated_at":"` + updated_at + `"}`
	w.Write([]byte(json))
	log.Printf("get repo info : %v \n", json)
}
