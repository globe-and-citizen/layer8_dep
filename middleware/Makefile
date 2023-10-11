.PHONY: help
.DEFAULT_GOAL := help
ARG :=

define base64_encode
	mkdir -p ./dist/`dirname $(1) | cut -d'/' -f3` && \
	base64 -w 0 $(1) | sed 's/^/"/' | sed 's/$$/"/' > ./dist/`basename $(1) | cut -d'.' -f1`.json && \
	echo "Encoded $(1) to ./dist/`basename $(1) | cut -d'.' -f1`.json";
endef

help: ## Show this help message
	@awk 'BEGIN {FS = ":.*##"; printf "Usage: make \033[36m<target>\033[0m\n\n"} /^[a-zA-Z_-]+:.*?##/ { printf "\033[36m%-10s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

build: ## Build WASM middleware for the current version
	@GOOS=js GOARCH=wasm go build -o ./bin/middleware.wasm ./main.go
	@echo "Built ./bin/middleware.wasm. Encoding..."
	@make encode ARG=./bin/middleware.wasm

encode: ## Encode the file specified by ARG or all files in ./bin if no ARG is specified
	@if [ -z "$(ARG)" ]; then \
		for file in `find ./bin -type f`; do \
			$(call base64_encode,$$file) \
		done \
	else \
		$(call base64_encode,$(ARG)) \
	fi
