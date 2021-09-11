echo "start vagrant provioning..."
sudo apt-get update

echo "configure openshift user and group"
sudo adduser vagrant openshift
newgrp openshift

echo "install openshift tools"
sudo apt install zfsutils-linux

echo "setup openshift manager"
openshift init --preseed

echo "install curl pkg..."
sudo apt-get install -y curl zfsutils-linux

echo "install golang pkg"
sudo add-apt-repository ppa:longsleep/golang-backports
sudo apt update -y
sudo apt install -y golang-go 

echo "Install dlv pkg"
 git clone https://github.com/go-delve/delve.git $GOPATH/src/github.com/go-delve/delve
 cd $GOPATH/src/github.com/go-delve/delve
 make install

### export dlv bin path
export PATH=$PATH:/home/vagrant/go/bin

echo "Finished provisioning."
