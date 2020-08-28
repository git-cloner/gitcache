package main

import (
	"io"
	"log"
	"net/http"
	"strings"
)

func rinetGitRequest(w http.ResponseWriter, r *http.Request, url string) {
	log.Printf("redirect to github.com : %v,%v\n", url, r.Method)
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
		panic("redirect to github.com to be an http.Flusher")
	}
	p := make([]byte, 204800)
	for {
		n_read, err := resp.Body.Read(p)
		//log.Printf("clone from github.com direct : %v,%v\n", url, n_read)
		if n_read > 0 {
			n_write, err := w.Write(p[:n_read])
			if err != nil {
				panic("redirect to github.com with pipe error:" + err.Error())
			}
			if n_read != n_write {
				panic("redirect to github.com with pipe failed to write data")
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
		return realurl
	} else {
		return ""
	}
}

func RequestHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "gitcache/download") {
			DownloadFile(w, r)
			return
		}
		var url = preProcUrl(r.URL.RequestURI())
		if url == "" {
			log.Printf("unknown token : %v\n", r.URL.RequestURI())
			w.WriteHeader(404)
			return
		}
		url = "https://" + url
		log.Printf("client send git request: %s\n", url)
		hdrNocache(w)
		//redirect to github.com clone
		rinetGitRequest(w, r, url)
		return
	}
}
