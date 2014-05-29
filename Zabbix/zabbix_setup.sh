#!/bin/bash

# Installing Zabbix agent on supervised node:
# This script is intended to be run on the "agent" node
# Must be run as user root

# Getting some info
node=`uname -n`
# ip=`ifconfig eth0 | grep "inet addr"|cut -f2 -d:|cut -f1 -dB`

# Create user environment
groupadd zabbix
useradd -g zabbix zabbix
mkdir /home/zabbix
chown zabbix:zabbix /home/zabbix/

# copy bin-file
cd /usr/local/sbin/
scp root@10.0.0.3:/usr/local/sbin/zabbix_a* .

# copy and change configuration file
cd ../etc/
scp root@10.0.0.3:/usr/local/etc/zabbix_a* .
sed "s/Hostname=Zabbix server/Hostname=$node/g" zabbix_agentd.conf >temp
sed "s/Server=127.0.0.1/Server=10.0.0.3/g" temp >temp2
sed "s/ServerActive=127.0.0.1/ServerActive=10.0.0.3/g" temp2 >temp3
mv temp3 zabbix_agentd.conf

# Open firewall
ufw allow 10050
ufw allow 10051

# Setup autostart
cd /etc/init.d/
scp root@10.0.0.3:/etc/init.d/zabbix_a* .
cd ../rc0.d
ln -s ../init.d/zabbix_agent K20zabbix-agent
cd ../rc1.d/
ln -s ../init.d/zabbix_agent K20zabbix-agent
cd ../rc2.d/
ln -s ../init.d/zabbix_agent S20zabbix-agent
cd ../rc3.d/
ln -s ../init.d/zabbix_agent S20zabbix-agent
cd ../rc4.d/
ln -s ../init.d/zabbix_agent S20zabbix-agent
cd ../rc5.d/
ln -s ../init.d/zabbix_agent S20zabbix-agent
cd ../rc6.d/
ln -s ../init.d/zabbix_agent K20zabbix-agent
