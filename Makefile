SHELL := /bin/bash

GOCMD=go
MOVESANDBOX=mv ~/vms/openshift-ordealopenshift-ordeal ~/vms-local/openshift-ordeal
GOPACKR=$(GOCMD) get -d github.com/gobuffalo/packr/packr && ${GOPATH}/bin/packr
GOMOD=$(GOCMD) mod
GOMOCKS=$(GOCMD) generate ./...
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
BINARY_NAME=openshift-ordeal
GOCOPY=cp openshift-ordeal ~/vagrant_file/.

all:test lint build

fmt:
	$(GOCMD) fmt ./...
lint:
	./scripts/lint.sh
tidy:
	$(GOMOD) tidy -v
test:
	$(GOCMD) get -d github.com/golang/mock/mockgen@v1.6.0
	$(GOCMD) install -v github.com/golang/mock/mockgen && export PATH=$GOPATH/bin:$PATH;
	$(GOMOCKS)
	$(GOTEST) ./... -coverprofile coverage.md fmt
	$(GOCMD) tool cover -html=coverage.md -o coverage.html
	$(GOCMD) tool cover  -func coverage.md
build:
	$(GOPACKR)
	export PATH=$GOPATH/bin:$PATH;
	export PATH=$PATH:/home/vagrant/go/bin
	export PATH=$PATH:/home/root/go/bin
	GOOS=linux GOARCH=amd64 $(GOBUILD) -v ./cmd/openshift-ordeal;
build_local:
	packr
	export PATH=$GOPATH/bin:$PATH;
	export PATH=$PATH:/home/vagrant/go/bin
	export PATH=$PATH:/home/root/go/bin
	$(GOBUILD) ./cmd/openshift-ordeal;
install:build_travis
	cp $(BINARY_NAME) $(GOPATH)/bin/$(BINARY_NAME)
test_build_travis:
	$(GOCMD) get -d github.com/golang/mock/mockgen@v1.6.0
	$(GOCMD) install -v github.com/golang/mock/mockgen && export PATH=$GOPATH/bin:$PATH;
	$(GOMOCKS)
	$(GOTEST) -short ./...  -coverprofile coverage.md fmt
	$(GOCMD) tool cover -html=coverage.md -o coverage.html
	GOOS=linux GOARCH=amd64 $(GOBUILD) -v ./cmd/openshift-ordeal;
build_travis:
	$(GOPACKR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) -v ./cmd/openshift-ordeal;
build_remote:
	$(GOPACKR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) -v ./cmd/openshift-ordeal
	mv openshift-ordeal ~/boxes/basic_box/openshift-ordeal

build_docker_local:
	docker build -t chenkeinan/openshift-ordeal:3 .
	docker push chenkeinan/openshift-ordeal:3
dlv:
	dlv --listen=:2345 --headless=true --api-version=2 --accept-multiclient exec ./openshift-ordeal
build_beb:
	$(GOPACKR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) -v -gcflags='-N -l' cmd/openshift-ordeal/openshift-ordeal.go
	scripts/deb.sh
.PHONY: all build install test
