# Copyright 2019 Hayo van Loon
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

PROJECT_NAME := bobsknobshop
MODULE_NAME := commentcentre
VERSION := v1

# Docker-related
IMAGE_NAME := $(PROJECT_NAME)_$(MODULE_NAME)_$(VERSION)
TAG := latest

LOCAL_PORT := 9010
NODE_PORT := 30100

.PHONY:

clean:
	go clean

protoc:
	@echo Go App, skipped

test:
	@echo TODO

build:
	docker build -t $(IMAGE_NAME) .

build-all: build

deploy:
	kubectl apply -f deploy-service.yaml

run:
	go run server.go \
		-port=$(LOCAL_PORT) \
		-classy-host=localhost \
		-classy-port=9000

smoke-test:
	go run client/client.go \
		-port=$(LOCAL_PORT) \
		-file=$(shell echo $$PWD)/../../example/comments.csv

docker-run:
	docker run --network="host" $(IMAGE_NAME) \
		/usr/local/bin/app \
		-port=$(LOCAL_PORT)

test-minikube:
	go run client/client.go \
		-host=$(shell minikube ip) \
		-port=$(NODE_PORT) \
		-file=$(shell echo $$PWD)/../../example/comments.csv

push-gcr:
	docker tag $(IMAGE_NAME) gcr.io/$(PROJECT_ID)/$(IMAGE_NAME):$(TAG)
	docker push gcr.io/$(PROJECT_ID)/$(IMAGE_NAME)

update-deps:
	go get -u ./...
