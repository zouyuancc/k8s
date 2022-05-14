FROM ubuntu:18.04
RUN sed -i 's/archive.ubuntu.com/mirrors.ustc.edu.cn/g' /etc/apt/sources.list && \
apt-get update && \
  DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
  build-essential \
  apt-utils \
  ca-certificates \
  wget \
  git \
  vim \
  libssl-dev \
  curl \
  unzip \
  unrar \
  tmux \
  openssh-server \
  cmake \
  gdb \
  mysql-server \
  mysql-client && \
  usermod -d /var/lib/mysql/ mysql && \
  mkdir ~/mytmp && \
  echo 'root:123456' | chpasswd && \
  sed -i 's/#Port 22/Port 22/g' /etc/ssh/sshd_config && \
  sed -i 's/#PubkeyAuthentication yes/PubkeyAuthentication yes/g' /etc/ssh/sshd_config && \
  sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin yes/g' /etc/ssh/sshd_config && \
  sed -i 's/#PasswordAuthentication yes/PasswordAuthentication yes/g' /etc/ssh/sshd_config && \
  echo 'if [ $(ps -aux | grep sshd | wc -l) -eq 1 ];then service ssh start; fi' >> ~/.bashrc && \
  echo 'if [ $(ps -aux | grep mysqld | wc -l) -eq 1 ];then service mysql start; fi' >> ~/.bashrc
EXPOSE 3306
EXPOSE 22
EXPOSE 8888
RUN ldconfig && apt-get clean && apt-get autoclean && \
  rm -rf /var/cache/debconf/* \
  /var/cache/apt \
  /var/lib/apt/periodic \
  /var/lib/apt/lists/* \
  /tmp/* \
  ~/mytmp

## build an image
# docker build -t shiyf/go .
## run a container
# docker run -dit --rm -v /data1/shiyf/Dockerfile/shiyf-Golang/Golang:/root/data -P -p 1234:22 -e MYSQL_ROOT_PASSWORD=123456 --name shiyf shiyf/go /bin/bash