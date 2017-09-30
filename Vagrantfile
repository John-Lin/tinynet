# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|
  config.vm.box = "ubuntu/xenial64"
  config.vm.hostname = 'dev'
  config.vm.network "public_network"

  config.vm.provision "shell", privileged: false, inline: <<-SHELL
    set -e -x -u
    sudo apt-get update
    sudo apt-get install -y vim git build-essential openvswitch-switch bridge-utils

    # Install Golang
    wget --quiet https://storage.googleapis.com/golang/go1.9.linux-amd64.tar.gz
    sudo tar -zxf go1.9.linux-amd64.tar.gz -C /usr/local/

    echo 'export GOROOT=/usr/local/go' >> /home/ubuntu/.bashrc
    echo 'export GOPATH=$HOME/go' >> /home/ubuntu/.bashrc
    echo 'export PATH=$PATH:$GOROOT/bin:$GOPATH/bin' >> /home/ubuntu/.bashrc
    source /home/ubuntu/.bashrc

    mkdir -p /home/ubuntu/go/src

    rm -rf /home/ubuntu/go1.9.linux-amd64.tar.gz
  SHELL

  config.vm.provider :virtualbox do |v|
    v.customize ["modifyvm", :id, "--cpus", 2]
    # enable this when hosts up to 4 
    # v.customize ["modifyvm", :id, "--memory", 4096]
    v.customize ['modifyvm', :id, '--nicpromisc1', 'allow-all']
  end
end
