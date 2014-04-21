FROM ubuntu

RUN apt-key adv --keyserver keyserver.ubuntu.com --recv-keys B97B0AFCAA1A47F044F244A07FCC7D46ACCC4CF8
RUN echo "deb http://apt.postgresql.org/pub/repos/apt/ precise-pgdg main" > /etc/apt/sources.list.d/pgdg.list

RUN apt-get update
RUN apt-get -y -q install python-software-properties software-properties-common
RUN apt-get -y -q install postgresql-9.3 postgresql-client-9.3 postgresql-contrib-9.3
RUN apt-get update
RUN add-apt-repository -y ppa:chris-lea/node.js
RUN apt-get update
RUN apt-get install -y nodejs git tree

# Append to $PATH variable.
RUN echo '\n# Node.js\nexport PATH="node_modules/.bin:$PATH"' >> /root/.bash_profile

RUN echo '#! /bin/bash' > /home/launch.sh
RUN echo "timeout 20 /usr/lib/postgresql/9.3/bin/postgres -D /var/lib/postgresql/9.3/main -c config_file=/etc/postgresql/9.3/main/postgresql.conf  &" >> /home/launch.sh
RUN echo "timeout 2 node /home/M7011E/createTables ">> /home/launch.sh
#ADD launch.sh /script/launch.sh
RUN chmod 777 /home/launch.sh


RUN cd /home/ ; git clone https://github.com/AmelieA/M7011E.git
RUN cd /home/M7011E ; npm install

USER postgres

RUN    /etc/init.d/postgresql start &&\
    psql --command "CREATE USER nodetest WITH SUPERUSER PASSWORD 'pika';" &&\
    createdb -O nodetest dbtest
RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/9.3/main/pg_hba.conf

RUN echo "listen_addresses='*'" >> /etc/postgresql/9.3/main/postgresql.conf

# Expose the PostgreSQL port
#EXPOSE 5432
EXPOSE 8080

# Add VOLUMEs to allow backup of config, logs and databases
VOLUME	["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

# Set the default command to run when starting the container
#CMD ["/usr/lib/postgresql/9.3/bin/postgres", "-D", "/var/lib/postgresql/9.3/main", "-c", "config_file=/etc/postgresql/9.3/main/postgresql.conf"]
#CMD ["node","/root/M7011E/createTables"]
#CMD ["tree","-L","2","/root/"]
CMD ["./home/launch.sh"]
