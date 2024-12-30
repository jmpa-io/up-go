
# Default PROJECT, if not given by another Makefile.
ifndef PROJECT
PROJECT=up-go
endif

# Targets.
example: binary-go-example ## Build the example binary.
run: build-go-example ## Builds & runs the example binary.

---: ## ---

# Includes the common Makefile.
# NOTE: this recursively goes back and finds the `.git` directory and assumes
# this is the root of the project. This could have issues when this assumtion
# is incorrect.
include $(shell while [[ ! -d .git ]]; do cd ..; done; pwd)/Makefile.common.mk

