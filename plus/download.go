package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

func validRemoteFile(url string) bool {
	resp, err := http.Head(url)
	if err != nil {
		log.Println(err)
		return false
	}
	if resp.StatusCode != 404 {
		return true
	} else {
		return false
	}
}

func ParseToken(tokenString string) (jwt.MapClaims, error) {
	SECRETKEY := "243223ffslsfsldfl412fdsfsdf"
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(SECRETKEY), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
func ValidUser(token string) bool {
	_, error := ParseToken(token)
	log.Printf("validUser error message : %v \n", error)
	return error == nil
}

func DownloadFile(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	var token = ""
	if len(params["token"]) > 0 {
		token = params["token"][0]
	}
	url := "https:" + strings.Replace(strings.Replace(r.URL.RequestURI(), "gitcache/download", "", -1), "?token="+token, "", -1)
	if !strings.Contains(url, "://github.com/") && !strings.Contains(url, "://raw.githubusercontent.com/") {
		log.Printf("invalid url")
		w.WriteHeader(401)
		return
	}
	if !ValidUser(token) {
		log.Printf("invalid user info")
		w.WriteHeader(403)
		return
	}
	if !validRemoteFile(url) {
		log.Printf("redirect to : %v 404\n", url)
		w.WriteHeader(404)
		return
	}
	log.Printf("redirect to : %v\n", url)
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
		panic("redirect to " + url + "  be an http.Flusher")
	}
	p := make([]byte, 20480)
	for {
		n_read, err := resp.Body.Read(p)
		if n_read > 0 {
			n_write, err := w.Write(p[:n_read])
			if err != nil {
				panic("redirect to " + url + " with pipe error:" + err.Error())
			}
			if n_read != n_write {
				panic("redirect to" + url + " with pipe failed to write data")
			}
			flusher.Flush()
		}
		if err == io.EOF {
			break
		}
	}
}
