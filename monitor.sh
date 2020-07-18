#!/bin/bash
COUNT=`netstat -ano|grep 5000 | grep LISTEN | wc -l`
if [ "$COUNT" -eq 0 ];
then
  echo "gitcache servicie not exist" 
  cd ~
  nohup ~/gitcache/gitcache -b /home/gitclone/repos/gitcache > gitcache.log 2>&1 &
fi
