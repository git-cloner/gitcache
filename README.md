# Gitcache
github.com clone cache

use git http  protocol to proxy git clone.

When the local cache does not exist, the clone request is redirected to github.com, and the mirror is created at same time(delay 10 seconds), and the next time it is cloned, then clone from the local mirror .

## build

install go environment,and

```shell
#clone
git clone https://github.com/git-cloner/gitcache
cd gitcache
#linux
export GO111MODULE=on
export GOPROXY=https://goproxy.io
#windows
set GO111MODULE=on
set GOPROXY=https://goproxy.io
#build
go build
```

## run

```shell
# -b git cahce base path
#linux
./gitcache  -b /var/gitcache
#windows
gitcache -b d:\temp
```

 

## use

git clone http://127.0.0.1:5000/github.com/git-cloner/gitcache

## homepage

and please try https://gitclone.com/ 

## client

you can use cgit client. https://github.com/git-cloner/gitcache/releases/download/v0.1/cgit-release.zip

cgit clone https://github.com/git-cloner/gitcache