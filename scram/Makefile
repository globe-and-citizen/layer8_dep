.PHONY: help
.DEFAULT_GOAL := help
ARG :=

help: ## Show this help message
	@awk 'BEGIN {FS = ":.*##"; printf "Usage: make \033[36m<target>\033[0m\n\n"} /^[a-zA-Z_-]+:.*?##/ { printf "\033[36m%-10s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

build: ## Build WASM 
	@GOOS=js GOARCH=wasm go build -o ../resource_server/frontend/public/bin/scram.wasm ./main.go
	@echo "Built ../resource_server/frontend/public/bin/scram.wasm"
