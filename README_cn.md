# Gitcache
[英文说明](https://github.com/git-cloner/gitcache/blob/master/README.md)

github.com clone 缓存，使用git的http协议代理git clone操作，当本地镜像（缓存）未建立前，clone操作被重定向到github.com，镜像会延时10秒开始创建，缓存建立后，下次clone（其他开发者）时就会利用到本地缓存，每晚自动从github.com更新镜像。

## 安装golang环境（linux）

```shell
#download golang,用登录用户，不要用sudo或root用户
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

如果用windows,从 https://golang.google.cn/dl/ 下载windows安装包安装。

## 编译

clone代码，然后设置环境变量支持go的 module模式。

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

## 运行

```shell
# -b git cahce base path
#linux
./gitcache  -b /var/gitcache
#windows
gitcache -b d:\temp
#重要提醒
用HTTPS git remote update 非常困难,好消息是可以用ssh,比如：
gitcache -ssh 1 -b /var/gitcache
使用之前要配好ssh,并执行以下命令将私钥加到高速缓存
eval $(ssh-agent -s)
ssh-add ~/.ssh/id_rsa
#数据库支持
如果要使用数据库存储缓存的一些信息（非必需）,需要设置环境变量
export MYSQL_DSN=dbuser:password@tcp(IP:3306)/dbname
#删除空文件夹
cd ~/repos/gitcache/github.com
find  -maxdepth 1 -type  d -empty | xargs -n 1 rm -fr
```

 

## 使用

简单修改url即可。

git clone http://127.0.0.1:5000/github.com/git-cloner/gitcache

## 利用gitcache技术建立的网站

https://gitclone.com/ 

## 客户端支持

从  https://github.com/git-cloner/gitcache/releases/download/v0.1/cgit-release.zip 下载，只要把git换成cgit即可加速，非常简单。   
或者通过命令行将`cgit`客户端程序安装到`GOBIN`目录下
```sh
go install github.com/git-cloner/gitcache/cgit@latest
```

```sh
cgit clone https://github.com/git-cloner/gitcache
```

## 区块链技术

gitcache的分布式缓存协调共享机制，使用了codechain技术(基于tendermint构建)

[https://github.com/little51/codechain](https://github.com/little51/codechain)
