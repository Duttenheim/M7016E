#!/bin/bash

ping -c 10 google.com >> google.txt &
ping -c 5 wikepedia.org >> wiki.txt
ping -c 5 google.com &
ping -c 5 google.fr

sleep 10

tree -L 2 /root/

cat google.txt
cat wiki.txt
