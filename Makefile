# This is the Makefile for the PLACEHOLDER project.
# Its features are described in the (project) Readme.

MAKE := /usr/bin/make

# General Recipes
dev-env:
	# add specific recipes here as needed / convenient

test: dev-env
	# add specific recipes here as needed / convenient

clean:
	make -C classy clean
	make -C peddler clean
	make -C truth clean

protoc:
	# add specific recipes here as needed / convenient
	make -C proto protoc-go

# Module-Specific Recipes
