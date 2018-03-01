#!/bin/bash -v

sudo yum install -y gcc python27 python27-devel python27-pip libffi-devel openssl-devel git httpd

sudo /usr/bin/pip-2.7 install --upgrade ansible

sudo /usr/bin/pip-2.7install ansible-lint

sudo /usr/bin/pip-2.7 install ansible-review 

sudo /usr/bin/pip-2.7 install bigsuds

sudo /usr/bin/pip-2.7 install f5-sdk

cd /home/ec2-user

/usr/bin/wget http://mirrors.jenkins.io/war-stable/latest/jenkins.war

sudo service httpd start

sudo yum install -y tomcat8

cd /usr/share/tomcat8/webapps

sudo /usr/bin/wget https://storage.googleapis.com/google-code-archive-downloads/v2/code.google.com/bodgeit/bodgeit.1.4.0.zip

sudo unzip bodgeit.1.4.0.zip

sudo service tomcat8 start

