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
