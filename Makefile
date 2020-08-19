PROJECTNAME=$(shell basename "$(PWD)")
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin

define RECOMMENDED_TAG=
$(git tag -l | sort -V | tail -n 1 | sed 's#\.# #g' | awk '{ print $1"."$2 + 1".0" }')
endef

all: install build docker-build docker-push

install:
	go mod download

dev:
	go build -o $(GOBIN)/app ./main.go && FPM_DEVICE=foodevice $(GOBIN)/app

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(GOBIN)/app ./main.go
	tar -zcvf app.tar.gz $(GOBIN)/app

build-prod:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -tags prod -o $(GOBIN)/app ./main.go