# openshift-vagrantfile

vagrant file to be used for openshift associated  programs developments, file include :
- buntu/bionic64
- openshift cluster 
- dlv for remote debug

### Quick Start

```
 git clone git@github.com:chen-keinan/openshift-vagrantfile.git
 cd openshift-vagrantfile
 vagrant up

```

### Compile binary with debug params
```
GOOS=linux GOARCH=amd64 go build -v -gcflags='-N -l' demo.go
```
### Run debug on remote machine
```
dlv --listen=:2345 --headless=true --api-version=2 --accept-multiclient exec ./demo
```

### Tear down
```
 vagrant destroy
