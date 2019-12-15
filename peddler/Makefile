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

# Generic Makefile
# Should contain superset of recipes accessible from project root.
#
# The type of module determines the behaviour of the recipe, i.e.: the 'run'
# recipe could just as well start a local appengine/flask as running a pipeline
# online (incurring costs). Check it before you run it.
#

CURRENT_VERSION := v1

dev-env:
	@$(MAKE) -C

clean:
	@$(MAKE) -C $(CURRENT_VERSION) clean

protoc:
	@$(MAKE) -C $(CURRENT_VERSION) protoc

build:
	@$(MAKE) -C $(CURRENT_VERSION) build

test:
	@$(MAKE) -C $(CURRENT_VERSION) test

dist:
	@$(MAKE) -C $(CURRENT_VERSION) dist

install:
	@$(MAKE) -C $(CURRENT_VERSION) dist

deploy:
	@$(MAKE) -C $(CURRENT_VERSION) deploy

run:
	@$(MAKE) -C $(CURRENT_VERSION) run

docker-run:
	@$(MAKE) -C $(CURRENT_VERSION) docker-run

data-local:
	@$(MAKE) -C $(CURRENT_VERSION) data-local

data-cloud:
	@$(MAKE) -C $(CURRENT_VERSION) data-online

train-local:
	@$(MAKE) -C $(CURRENT_VERSION) train-local

train-cloud:
	@$(MAKE) -C $(CURRENT_VERSION) train-online

hyper:
	@$(MAKE) -C $(CURRENT_VERSION) hyper

tensorboard:
	@$(MAKE) -C $(CURRENT_VERSION) tensorboard

tensorboard-local:
	@$(MAKE) -C $(CURRENT_VERSION) tensorboard

update-deps:
	@$(MAKE) -C $(CURRENT_VERSION) update-deps

smoke-test:
	@$(MAKE) -C $(CURRENT_VERSION) smoke-test

build-all:
	@$(MAKE) -C $(CURRENT_VERSION) build-all

deploy:
	@$(MAKE) -C $(CURRENT_VERSION) deploy
