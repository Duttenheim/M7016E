#!/bin/bash
# forked from https://gist.github.com/jpetazzo/5494158


DATADIR=/var/lib/postgresql/9.3/main
BINDIR=/usr/lib/postgresql/9.3/bin

# test if DATADIR is existent
#if [ ! -d $DATADIR ]; then
  echo "Creating Postgres data at $DATADIR"
  mkdir -p $DATADIR
#fi

echo 'host all all 0.0.0.0/0 md5' >> /etc/postgresql/9.3/main/pg_hba.conf

# test if DATADIR has content
#if [ ! "$(ls -A $DATADIR)" ]; then
  echo "Initializing Postgres Database at $DATADIR"
  chown -R postgres $DATADIR
  su postgres sh -c "$BINDIR/initdb $DATADIR"
  su postgres sh -c "$BINDIR/postgres --single  -D $DATADIR  -c config_file=/etc/postgresql/9.1/main/postgresql.conf" <<< "CREATE USER nodetest WITH SUPERUSER PASSWORD 'pika';"
  psql -U postgres -d dbtest -c "
CREATE TABLE Users(firstName VARCHAR(255), lastName VARCHAR(255), googleID VARCHAR(255) PRIMARY KEY, banned BOOL);"
  psql -U postgres -d dbtest -c "
CREATE TABLE Locations(location VARCHAR(255) PRIMARY KEY, x REAL NOT NULL, y REAL NOT NULL);"
  psql -U postgres -d dbtest -c "
CREATE TABLE Comments(location VARCHAR(255) REFERENCES Locations(location), text VARCHAR(255), login VARCHAR(255), date VARCHAR(255));"
  psql -U postgres -d dbtest -c "
CREATE TABLE Images(location VARCHAR(255) REFERENCES Locations(location), img_name VARCHAR(255), login VARCHAR(255), date VARCHAR(255));"
#fi

echo "ready to launch"
#cat /etc/postgresql/9.3/main/postgresql.conf

#su postgres sh -c "$BINDIR/postgres           -D $DATADIR  -c config_file=/etc/postgresql/9.3/main/postgresql.conf  -c listen_addresses=*" 
