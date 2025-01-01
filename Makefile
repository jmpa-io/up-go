
# Default PROJECT, if not given by another Makefile.
ifndef PROJECT
PROJECT=up-go
endif

# Targets.
accounts: binary-go-accounts ## Build the 'accounts' binary.
ping: binary-go-ping ## Build the 'ping' binary.
tags: binary-go-tags ## Build the 'tags' binary.
transactions: binary-go-transactions ## Build the `transactions` binary.
run: accounts ping tags transactions
PHONY += run

---: ## ---

# Includes the common Makefile.
# NOTE: this recursively goes back and finds the `.git` directory and assumes
# this is the root of the project. This could have issues when this assumtion
# is incorrect.
include $(shell while [[ ! -d .git ]]; do cd ..; done; pwd)/Makefile.common.mk

