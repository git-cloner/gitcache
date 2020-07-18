#!/bin/bash
cd ~
cd gitcache
git pull
go build
cd ~
pkill gitcache
ps -ef|grep gitcache/gitcache
./monitor.sh
tail -f gitcache.log
