FROM dockerfile/nodejs
# based on ubuntu 12.10


#RUN sudo apt-get install tree
#RUN git clone https://github.com/AmelieA/M7011E.git
#RUN cd M7011E ; npm install

RUN apt-key adv --keyserver keyserver.ubuntu.com --recv-keys B97B0AFCAA1A47F044F244A07FCC7D46ACCC4CF8
RUN echo "deb http://apt.postgresql.org/pub/repos/apt/ precise-pgdg main" > /etc/apt/sources.list.d/pgdg.list
RUN apt-get update
RUN apt-get -y -q install python-software-properties software-properties-common
RUN apt-get -y -q install postgresql-9.3 postgresql-client-9.3 postgresql-contrib-9.3
#RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/9.3/main/pg_hba.conf

#ADD start.sh /root/start.sh
#RUN chmod 777 /root/start.sh

#USER postgres
#RUN    /etc/init.d/postgresql start &&\
#    psql --command "CREATE USER nodetest WITH SUPERUSER PASSWORD 'pika';" &&\
#    createdb -O nodetest dbtest
#EXPOSE 8080

ADD launch.sh /root/launch.sh
RUN chmod 777 /root/launch.sh
#CMD ["node","/root/M7011E/app"]
#CMD ["tree","-L","2","/root/"]
CMD ["/root/launch.sh"]
