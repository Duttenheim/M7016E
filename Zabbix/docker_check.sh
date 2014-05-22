#!/bin/bash

# Find docker container
container=$1
item=$2

dockerCnt=`sudo docker ps | grep $container | wc -l`
if [ $dockerCnt != 1 ]; then
  echo "-1"
  exit
fi

containerId=`sudo docker ps|grep $container|sed 's/^ *//'|cut -f 1 -d ' '`

pid=`sudo docker inspect $containerId|grep Pid|cut -f 2 -d ':'|cut -f 1 -d ','`

list=$pid
do_pid() {
  local ppid=$1
  for line in `ps -e -o pid,ppid|grep $ppid|sed 's/^ *//'|cut -f 1 -d ' '` ; do
    pid1=`echo $line | cut -f 1 -d ' '`
    pid2=`echo $line | cut -f 2 -d ' '`
    if [ "$pid1" != "$ppid" ] ; then
      list="$list $pid1"
      do_pid "$pid1"
    fi
  done
}

do_pid "$pid"

MEM=0.0
CPU=0.0

for pid in `echo $list` ; do
  mem=`ps -p $pid -o %mem | tail -1 | sed 's/^ *//' | cut -f 1 -d ' '`
  cpu=`ps -p $pid -o %cpu | tail -1 | sed 's/^ *//' | cut -f 1 -d ' '`
  second=`echo $cpu | cut -b 2`
  if [ "$second" != "C" ] ; then
    MEM=`echo "$MEM + $mem" | bc`
    CPU=`echo "$CPU + $cpu" | bc`
  fi
done

case "$item" in
  cpu)
    echo $CPU
    ;;
  mem)
    echo $MEM
    ;;
  *)
    ;;
esac
