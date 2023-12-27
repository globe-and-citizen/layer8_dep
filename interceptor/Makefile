include .env
export

.PHONY := help
.DEFAULT_GOAL := help
ARG :=

# ALTERED PATH
define base64_encode
	mkdir -p ../proxy/assets/cdn/interceptor/`dirname $(1) | cut -d'/' -f3` && \
	base64 -w 0 $(1) | sed 's/^/"/' | sed 's/$$/"/' > ../server/assets-v1/cdn/interceptor/`basename $(1) | cut -d'.' -f1`.json && \
	echo "Encoded $(1) to ../server/assets-v1/cdn/interceptor/`basename $(1) | cut -d'.' -f1`.json";
endef

help: ## Show this help message
	@awk 'BEGIN {FS = ":.*##"; printf "Usage: make \033[36m<target>\033[0m\n\n"} /^[a-zA-Z_-]+:.*?##/ { printf "\033[36m%-10s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

build: build_global build_local ## Build WASM Interceptor for global and local use

build_global: ## Build WASM Interceptor for global use
	@GOOS=js GOARCH=wasm go build -ldflags="-X main.Layer8Scheme=$(LAYER8_PROXY_SCHEME) -X main.Layer8Host=$(LAYER8_PROXY_DOMAIN) -X main.Layer8Port=$(LAYER8_PROXY_PORT)" -o ./bin/interceptor.wasm ./interceptor.go
	@echo "Built ./bin/interceptor.wasm. Encoding..."
	@make encode ARG=./bin/interceptor.wasm

build_local: ## Build WASM Interceptor for local use
	@GOOS=js GOARCH=wasm go build -ldflags="-X main.Layer8Scheme=$(LAYER8_PROXY_SCHEME_LOCAL) -X main.Layer8Host=$(LAYER8_PROXY_DOMAIN_LOCAL) -X main.Layer8Port=$(LAYER8_PROXY_PORT_LOCAL)" -o ./bin/interceptor__local.wasm ./interceptor.go
	@echo "Built ./bin/interceptor__local.wasm. Encoding..."
	@make encode ARG=./bin/interceptor__local.wasm

encode: ## Encode the file specified by ARG or all files in ./bin if no ARG is specified
	@if [ -z "$(ARG)" ]; then \
		for file in `find ./bin -type f`; do \
			$(call base64_encode,$$file) \
		done \
	else \
		$(call base64_encode,$(ARG)) \
	fi

build_local_cheating: ## Build the WASM interceptor without binary conversion
	@GOOS=js GOARCH=wasm go build -ldflags="-X main.Layer8Scheme=$(LAYER8_PROXY_SCHEME_LOCAL) -X main.Layer8Host=$(LAYER8_PROXY_DOMAIN_LOCAL) -X main.Layer8Port=$(LAYER8_PROXY_PORT_LOCAL)" -o ../proxy/assets/cdn/interceptor/interceptor__local.wasm ./interceptor.go
