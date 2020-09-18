package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type HttpParams struct {
	Repository string
	Gitservice string
	IsInfoReq  bool
}

func parseHttpParams(realurl string) HttpParams {
	u, err := url.Parse(realurl)
	if err != nil {
		panic(err)
	}
	str := strings.Split(u.Path, "/")
	if len(str) < 4 {
		panic("bad request params")
	}
	_Repository := str[1] + "/" + str[2] + "/" + str[3]
	var _Gitservice = strings.Replace(u.RawQuery, "service=", "", -1)
	if _Gitservice == "" {
		if (strings.Index(str[4], "git") != -1) && (strings.Index(str[4], "pack") != -1) {
			_Gitservice = str[4]
		}
	}
	_IsInfoReq := (str[4] == "info")
	var httpParams HttpParams = HttpParams{Repository: _Repository, Gitservice: _Gitservice, IsInfoReq: _IsInfoReq}
	return httpParams
}

func rinetGitRequest(w http.ResponseWriter, r *http.Request, url string) {
	log.Printf("redirect to : %v,%v\n", url, r.Method)
	client := &http.Client{}
	req, err := http.NewRequest(r.Method, url, r.Body)
	for k, v := range r.Header {
		for _, vv := range v {
			req.Header.Add(k, vv)
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	for k, v := range resp.Header {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}
	//
	flusher, ok := w.(http.Flusher)
	if !ok {
		panic("redirect to be an http.Flusher")
	}
	p := make([]byte, 20480)
	for {
		n_read, err := resp.Body.Read(p)
		//log.Printf("clone from github.com direct : %v,%v\n", url, n_read)
		if n_read > 0 {
			n_write, err := w.Write(p[:n_read])
			if err != nil {
				panic("redirect to with pipe error:" + err.Error())
			}
			if n_read != n_write {
				panic("redirect to with pipe failed to write data")
			}
			flusher.Flush()
		}
		if err == io.EOF {
			break
		}
	}
}

func hdrNocache(w http.ResponseWriter) {
	w.Header().Set("Expires", "Fri, 01 Jan 1980 00:00:00 GMT")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Cache-Control", "no-cache, max-age=0, must-revalidate")
}

func cors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("content-type", "application/json")
}

func RequestFromRemote(url string) *http.Response {
	client := &http.Client{}
	reqest, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	reqest.Header.Add("User-Agent", "git/")
	response, err1 := client.Do(reqest)
	if err1 != nil {
		panic(err1)
	}
	defer response.Body.Close()
	return response
}

func preProcUrl(url string) string {
	var realurl = ""
	var token = ""
	comma := strings.Index(url, ".github.com")
	if comma > 1 {
		token = url[1:comma]
		realurl = url[comma+1:]
	}
	log.Printf("real url: %s,token: %s\n", realurl, token)
	if ValidUser(token) {
		return "/" + realurl
	} else {
		return ""
	}
}

func imageExists(url string) bool {
	resp, err := http.Get("https://gitclone.com/gitcache/image/" + url)
	if err != nil {
		log.Println(err)
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func RequestHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "gitcache/download") {
			DownloadFile(w, r)
			return
		}
		if strings.Contains(r.URL.Path, "gitcache/star") {
			GetRepoStar(w, r)
			return
		}
		if strings.Contains(r.URL.Path, "/favicon.ico") {
			return
		}
		var realurl = preProcUrl(r.URL.RequestURI())
		if realurl == "" {
			log.Printf("unknown token : %v\n", r.URL.RequestURI())
			w.WriteHeader(404)
			return
		}
		var httpParams HttpParams = parseHttpParams(realurl)
		repos := httpParams.Repository
		if strings.HasSuffix(repos, ".git") {
			repos = repos[:len(repos)-4]
		}
		useImage := imageExists(repos)
		log.Printf("check image exists: %s %v\n", repos, useImage)
		realurl = "https:/" + realurl
		//redirect to github.com or cache to clone
		if useImage {
			realurl = strings.Replace(realurl, "https://github.com", "http://gitclonecache.com:22003/github.com", -1)
		}
		log.Printf("client send : %s\n", realurl)
		hdrNocache(w)
		rinetGitRequest(w, r, realurl)
		return
	}
}
