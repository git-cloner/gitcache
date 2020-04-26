# gitcache
git clone cache

## build

install go lang environment,and

```shell
export GO111MODULE=on
export GOPROXY=https://goproxy.io
go build
```

## run

nohup ./gitcache  -b /var/gitcache > gitcache.log 2>&1 &

## use

git clone http://127.0.0.1:5000/github.com/yourrepository