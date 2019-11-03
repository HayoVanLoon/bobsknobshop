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

# General Recipes
dev-env:
	# add specific recipes here as needed / convenient

test: dev-env
	# add specific recipes here as needed / convenient

clean: clean-truth
	$(MAKE) -C classy clean
	$(MAKE) -C peddler clean

protoc:
	# add specific recipes here as needed / convenient
	$(MAKE) -C proto protoc-go

update-deps: update-deps-classy update-deps-peddler update-deps-truth


# Module-Specific Recipes

# Classy
run-classy:
	$(MAKE) -C classy run

update-deps-classy:
	@$(MAKE) -C classy update-deps


# Peddler
run-peddler:
	$(MAKE) -C peddler run

update-deps-peddler:
	@$(MAKE) -C peddler update-deps

smoke-test-peddler:
	@$(MAKE) -C peddler smoke-test


# Truth
clean-truth:
	$(MAKE) -C truth clean

run-truth:
	$(MAKE) -C truth run

update-deps-truth:
	@$(MAKE) -C truth update-deps

smoke-test-truth:
	@$(MAKE) -C truth smoke-test
