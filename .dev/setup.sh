echo "start vagrant provioning..."
sudo apt-get update -y
sudo apt-get install docker.io -y

echo "configure openshift user and group"
systemctl start docker
systemctl enable docker
systemctl status docker

echo "download openshift origin"
wget https://github.com/openshift/origin/releases/download/v3.11.0/openshift-origin-client-tools-v3.11.0-0cbc58b-linux-64bit.tar.gz

echo "extract openshift origin"
tar -xvzf openshift-origin-client-tools-v3.11.0-0cbc58b-linux-64bit.tar.gz
cd openshift-origin-client-tools-v3.11.0-0cbc58b-linux-64bit
cp oc kubectl /usr/local/bin/
oc version

echo "{\"insecure-registries\" : [ \"172.30.0.0/16\" ]}" > /etc/docker/daemon.json
systemctl restart docker
echo YKyk426144 > my_password.txt
cat ~/my_password.txt | docker login --username chenkeinan --password-stdin

oc cluster up --public-hostname=172.30.1.5
oc login -u system:admin
oc project default
oc status

oc new-project dev --display-name="Project - Dev" --description="My Project"
oc project my-project
oc status

oc tag --source=docker openshift/deployment-example:v2 deployment-example:latest
oc new-app deployment-example:latest
oc status

oc expose service/deployment-example

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
