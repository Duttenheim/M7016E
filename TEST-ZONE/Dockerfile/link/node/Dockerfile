FROM ubuntu

RUN apt-key adv --keyserver keyserver.ubuntu.com --recv-keys B97B0AFCAA1A47F044F244A07FCC7D46ACCC4CF8
RUN echo "deb http://apt.postgresql.org/pub/repos/apt/ precise-pgdg main" > /etc/apt/sources.list.d/pgdg.list

RUN apt-get update
RUN apt-get -y -q install python-software-properties software-properties-common
RUN add-apt-repository -y ppa:chris-lea/node.js
RUN apt-get update
RUN apt-get install -y nodejs git tree

# Append to $PATH variable.
RUN echo '\n# Node.js\nexport PATH="node_modules/.bin:$PATH"' >> /root/.bash_profile

RUN cd /home/ ; git clone https://github.com/AmelieA/M7011E.git
RUN cd /home/M7011E ; sed -i 's/require('pg').native/require('pg')/g' createTables.js ; npm install

EXPOSE 8080

RUN echo '#!/bin/bash' > /home/launch.sh
RUN echo "cd /home/M7011E" >> /home/launch.sh
RUN echo "npm install" >> /home/launch.sh
RUN echo "timeout 5 node /home/M7011E/createTables" >> /home/launch.sh
RUN echo "node /home/M7011E/app" >> /home/launch.sh
RUN chmod 777 /home/launch.sh

#CMD ["tree","-L","2","/script/"]
CMD ["node", "/home/M7011E/createTables"]
#CMD ["./home/launch.sh"]
