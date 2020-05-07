# Gitcache
github.com clone cache

use git http  protocol to proxy git clone.

When the local cache does not exist, the clone request is redirected to github.com, and the mirror is created at same time, and the next time it is cloned, then clone from the local mirror .

## build

install go environment,and

```shell
export GO111MODULE=on
export GOPROXY=https://goproxy.io
go build
```

## run

./gitcache  -b /var/gitcache 

## use

git clone http://127.0.0.1:5000/github.com/yourrepository

and please try https://gitclone.com/ 
and use cgit client. https://github.com/git-cloner/gitcache/releases/download/v0.1/cgit-release.zip
