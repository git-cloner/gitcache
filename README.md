# Gitcache
[中文说明](https://github.com/git-cloner/gitcache/blob/master/README_cn.md)

github.com clone cache

use git http  protocol to proxy git clone.

When the local cache does not exist, the clone request is redirected to github.com, and the mirror is created at same time(delay 10 seconds), and it is cloned the next time, then clone from the local mirror .

new support branch (git clone -b branchname) 

## install golang（linux）

```shell
#download golang,use normal user,don't use sudo
curl -O https://dl.google.com/go/go1.14.linux-amd64.tar.gz
tar -xvf go1.14.linux-amd64.tar.gz
## install golang
sudo rm -fr /usr/local/go
sudo mv go /usr/local
mkdir goApps
## init environment variable
echo "export GOPATH=~/goApps" >> ~/.profile
echo "export PATH=\$PATH:/usr/local/go/bin" >> ~/.profile
echo "export PATH=\$PATH:\$GOPATH/bin" >> ~/.profile
## apply variable
source ~/.profile
```

if you use windows,download from https://golang.google.cn/dl/ and install.

## build

```shell
#clone
git clone https://github.com/git-cloner/gitcache
cd gitcache
#linux
export GO111MODULE=on
export GOPROXY=https://goproxy.cn
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
#important hint
Using HTTPS git remote update is very difficult,so you can use ssh,such as
gitcache -ssh 1 -b /var/gitcache
before use ssh,please config ssh first,and execute
eval $(ssh-agent -s)
ssh-add ~/.ssh/id_rsa
#database support
if using database(mysql) support(Not necessary),please set up environment variable
export MYSQL_DSN=dbuser:password@tcp(IP:3306)/dbname
```

 

## usage

git clone http://127.0.0.1:5000/github.com/git-cloner/gitcache

git clone  http://127.0.0.1:5000/github.com/git-cloner/gitcache -b branch

## homepage

and please try https://gitclone.com/ 

## client

you can use cgit client. https://github.com/git-cloner/gitcache/releases/download/v0.1/cgit-release.zip

cgit clone https://github.com/git-cloner/gitcache

## block chain

use codechain(chain base on tendermint)

[https://github.com/little51/codechain](https://github.com/little51/codechain)
