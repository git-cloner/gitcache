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

func parseHttpParams(r *http.Request) HttpParams {
	u, err := url.Parse(r.RequestURI)
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

func rinetGitRequest(w http.ResponseWriter, r *http.Request) {
	url := "https:/" + r.URL.RequestURI()
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
	p := make([]byte, 20480)
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

func RequestFromRemote(r *http.Request) *http.Response {
	var url = "https:/" + r.URL.RequestURI()
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

func RequestHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "gitcache/download") {
			DownloadFile(w, r)
			return
		}
		log.Printf("client send git request: %s %s %s %s\n", r.RemoteAddr, r.Method, r.URL.Path, r.Proto)
		var httpParams HttpParams = parseHttpParams(r)
		log.Printf("git params: %+v\n", httpParams)
		if ((r.Method == "GET") && (httpParams.IsInfoReq)) || ((r.Method != "GET") && (!httpParams.IsInfoReq)) {
			log.Printf("client send git request: %s %v valid ok\n", r.Method, httpParams.IsInfoReq)
		} else {
			log.Printf("not supported request : %v %v\n", r.Method, httpParams.IsInfoReq)
			w.WriteHeader(500)
			return
		}
		//only support git-upload-pack because
		if httpParams.Gitservice != "git-upload-pack" {
			if httpParams.Gitservice == "git-receive-pack" {
				body := RequestFromRemote(r)
				w.WriteHeader(body.StatusCode)
				return
			} else {
				log.Printf("not supported request : %v %v\n", r.Method, httpParams.Gitservice)
				w.WriteHeader(500)
				return
			}
		}
		if httpParams.IsInfoReq {
			hdrNocache(w)
			//redirect to github.com clone
			rinetGitRequest(w, r)
		} else {
			hdrNocache(w)
			//redirect to github.com clone
			rinetGitRequest(w, r)
		}
		return
	}
}
