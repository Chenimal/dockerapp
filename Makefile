#!/bin/bash

build:
	@go build

build-docker:
	@env GOOS=linux GOARCH=amd64 go build -o dockerapp

build-image:
	@docker build --tag ceruntu ./files/docker/

run:
	@docker-compose up
