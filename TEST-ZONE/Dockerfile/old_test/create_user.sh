#!/bin/bash

sudo /etc/init.d/postgresql start
sudo psql -U postgres --command "CREATE USER dbtest WITH SUPERUSER PASSWORD 'pika';"
#sudo -u postgres createuser nodetest

