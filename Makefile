
# Default PROJECT, if not given by another Makefile.
ifndef PROJECT
PROJECT=up-go
endif

# Variables.
TOKEN := $(shell aws ssm get-parameter --name "/tokens/up" --query 'Parameter.Value' --output text --with-decryption)

# Targets.
accounts: binary-go-accounts ## Build the 'accounts' binary.
ping: binary-go-ping ## Build the 'ping' binary.
tags: binary-go-tags ## Build the 'tags' binary.
tracing: binary-go-tracing ## Build the `tracing` binary.
transactions: binary-go-transactions ## Build the `transactions` binary.
run: accounts ping tags transactions tracing

PHONY += accounts ping tags transactions tracing run

get-token: ## Retrieves the Up token from AWS SSM Parameter Store.
get-token:
	@echo "export UP_TOKEN=\"$(TOKEN)\""

---: ## ---

# Includes the common Makefile.
# NOTE: this recursively goes back and finds the `.git` directory and assumes
# this is the root of the project. This could have issues when this assumtion
# is incorrect.
include $(shell while [[ ! -d .git ]]; do cd ..; done; pwd)/Makefile.common.mk

