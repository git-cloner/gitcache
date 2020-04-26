package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
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

func execShell(cmd string, args []string) string {
	log.Printf("execute local git command : %v,%v\n", cmd, args)
	var command = exec.Command(cmd, args...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	var err = command.Start()
	if err != nil {
		return err.Error()
	}
	err = command.Wait()
	if err != nil {
		return err.Error()
	}
	return ""
}

func fetchMirrorFromRemote(remote string, local string, update string) string {
	localLockExist, _ := PathExists(local + "/shallow.lock")
	if localLockExist {
		return "valid local cache error : cache is locked,please wait"
	}
	//var args = "-C " + local + " remote set-url origin " + remote
	var err = execShell("git", []string{"-C", local, "remote", "set-url", "origin", remote})
	if err != "" {
		return err
	}
	//args = "-C " + local + " fetch "
	if update == "" {
		return execShell("git", []string{"-C", local, "fetch", "--depth=1"})
	} else {
		//return execShell("git", []string{"-C", local, "fetch", "--unshallow"})
		return execShell("git", []string{"-C", local, "remote", "update"})
	}
}

func cloneMirrorFromRemote(remote string, local string) string {
	return execShell("git", []string{"clone", "--depth=1", "--mirror", "--progress", remote, local})
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func validLocalCache(local string) bool {
	var err = execShell("git", []string{"-C", local, "remote"})
	if err == "" {
		return true
	} else {
		os.RemoveAll(local)
		return false
	}
}

func ifExistsLocalCache(local string) (bool, string) {
	localGitExist, _ := PathExists(local)
	//.git path exists
	if localGitExist {
		localLockExist, _ := PathExists(local + "/shallow.lock")
		if localLockExist {
			var err1 = "git cache is updating... ...,please wait"
			log.Println(err1)
			return true, err1
		} else {
			return validLocalCache(local), ""
		}
	} else {
		return false, ""
	}
}

func mirrorFromRemote(remote string, local string) string {
	var localExists, err = ifExistsLocalCache(local)
	if err != "" {
		return err
	}
	if localExists {
		log.Printf("valid local cache! .git path exists")
		return ""
	} else {
		log.Printf("valid local cache! .git path not exists")
		err = cloneMirrorFromRemote(remote, local)
		if err != "" {
			log.Printf("git command: clone from remote error : %s\n", err)
		}
		return err
	}
}

func execGitCommand(cmd string, version string, args []string) []byte {
	log.Printf("execute local git command : %v,%v\n", cmd, args)
	command := exec.Command(cmd, args...)
	if len(version) > 0 {
		command.Env = append(os.Environ(), fmt.Sprintf("GIT_PROTOCOL=%s", version))
	}
	out, err := command.Output()

	if err != nil {
		log.Printf("execGitCommand error: %v\n", err)
	}
	return out
}

func execShelldPipe(cmd string, args []string, w http.ResponseWriter, r *http.Request) {
	var command = exec.Command(cmd, args...)
	in, err := command.StdinPipe()
	if err != nil {
		log.Printf("execute shell with pipe error: %v\n", err)
	}
	stdout, err := command.StdoutPipe()
	if err != nil {
		log.Printf("execute shell with pipe error: %v\n", err)
	}
	err = command.Start()
	if err != nil {
		log.Printf("execute shell with pipe error: %v\n", err)
	}
	var reader io.ReadCloser
	switch r.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(r.Body)
		defer reader.Close()
	default:
		reader = r.Body
	}
	io.Copy(in, reader)
	in.Close()
	flusher, ok := w.(http.Flusher)
	if !ok {
		panic("execute shell with pipe expected http.ResponseWriter to be an http.Flusher")
	}
	p := make([]byte, 1024)
	for {
		n_read, err := stdout.Read(p)
		if err == io.EOF {
			break
		}
		n_write, err := w.Write(p[:n_read])
		if err != nil {
			log.Printf("execute shell with pipe error: %v\n", err)
			os.Exit(1)
		}
		if n_read != n_write {
			log.Printf("execute shell with pipe failed to write data: %d read, %d written\n", n_read, n_write)
			os.Exit(1)
		}
		flusher.Flush()
	}
	command.Wait()
}

func hdrNocache(w http.ResponseWriter) {
	w.Header().Set("Expires", "Fri, 01 Jan 1980 00:00:00 GMT")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Cache-Control", "no-cache, max-age=0, must-revalidate")
}

func RequestHandler(basedir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("client send git request: %s %s %s %s\n", r.RemoteAddr, r.Method, r.URL.Path, r.Proto)
		var httpParams HttpParams = parseHttpParams(r)
		log.Printf("git params: %+v\n", httpParams)
		if ((r.Method == "GET") && (httpParams.IsInfoReq)) || ((r.Method != "GET") && (!httpParams.IsInfoReq)) {
			log.Printf("client send git request: %s %v valid ok\n", r.Method, httpParams.IsInfoReq)
		} else {
			//w.WriteHeader(200)
			//return
			panic("not supported request")
		}
		//only support git-upload-pack because
		if httpParams.Gitservice != "git-upload-pack" {
			if httpParams.Gitservice == "git-receive-pack" {
				body := RequestFromRemote(r)
				w.WriteHeader(body.StatusCode)
				return
			} else {
				panic("not supported request " + httpParams.Gitservice)
			}
		}
		var remote = "https://" + httpParams.Repository
		var local = path.Join(basedir, httpParams.Repository)
		//fix go get command,repository not has .git suffix
		if !strings.HasSuffix(local, ".git") {
			local = local + ".git"
		}
		if httpParams.IsInfoReq {
			log.Printf("git mirror from remote : %s %s\n", remote, local)
			err := mirrorFromRemote(remote, local)
			if err != "" {
				w.WriteHeader(500)
				w.Write([]byte(err))
			} else {
				refs := execGitCommand(httpParams.Gitservice, "", []string{"--stateless-rpc", "--advertise-refs", local})
				hdrNocache(w)
				w.Header().Set("Content-Type", fmt.Sprintf("application/x-%s-advertisement", httpParams.Gitservice))
				w.WriteHeader(200)
				w.Write([]byte("001e# service=git-upload-pack\n0000"))
				w.Write(refs)
			}
		} else {
			hdrNocache(w)
			w.Header().Set("Content-Type", fmt.Sprintf("application/x-%s-result", httpParams.Gitservice))
			w.WriteHeader(200)
			execShelldPipe(httpParams.Gitservice, []string{"--stateless-rpc", local}, w, r)
		}
		return
	}
}
