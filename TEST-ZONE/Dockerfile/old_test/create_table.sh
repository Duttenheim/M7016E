#!/bin/sh

sudo su - postgres
psql -U postgres -d dbtest -c "
CREATE TABLE Users(firstName VARCHAR(255), lastName VARCHAR(255), googleID VARCHAR(255) PRIMARY KEY, banned BOOL);"
psql -U postgres -d dbtest -c "
CREATE TABLE Locations(location VARCHAR(255) PRIMARY KEY, x REAL NOT NULL, y REAL NOT NULL);"
psql -U postgres -d dbtest -c "
CREATE TABLE Comments(location VARCHAR(255) REFERENCES Locations(location), text VARCHAR(255), login VARCHAR(255), date VARCHAR(255));"
psql -U postgres -d dbtest -c "
CREATE TABLE Images(location VARCHAR(255) REFERENCES Locations(location), img_name VARCHAR(255), login VARCHAR(255), date VARCHAR(255));"


