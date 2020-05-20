# Gitcache
[英文说明](https://github.com/git-cloner/gitcache/blob/master/README.md)

github.com clone 缓存，使用git的http协议代理git clone操作，当本地镜像（缓存）未建立前，clone操作被重定向到github.com，镜像会延时10秒开始创建，缓存建立后，下次clone（其他开发者）时就会利用到本地缓存，每晚自动从github.com更新镜像。

## 编译

安装got语言环境，clone代码，然后设置环境变量支持go的 module模式。

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
```

 

## 使用

简单修改url即可。

git clone http://127.0.0.1:5000/github.com/git-cloner/gitcache

## 利用gitcache技术建立的网站

https://gitclone.com/ 

## 客户端支持

从  https://github.com/git-cloner/gitcache/releases/download/v0.1/cgit-release.zip 下载，只要把git换成cgit即可加速，非常简单。

cgit clone https://github.com/git-cloner/gitcache

## 区块链技术

gitcache的分布式缓存协调共享机制，使用了codechain技术(基于tendermint构建)

[codechain](https://github.com/little51/codechain)