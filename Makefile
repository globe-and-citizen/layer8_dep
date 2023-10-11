## NPM install
npm_install_all:
	cd sp_mock/frontend && npm install && cd ../backend && npm install 

## Interceptor Calls
build_interceptor: ## must do from a bash terminal
	@$(MAKE) -C ./interceptor build

## Build Middleware
build_middleware:
	@$(MAKE) -C ./middleware build

## Run Mock

run_frontend:
	cd sp_mock/frontend && npm run dev
	
run_backend:
	cd sp_mock/backend && npm run dev

## Run Proxy

run_proxy:
	cd proxy && go run main.go --server=proxy --port=5001

# Serve auth server
run_auth:
	cd proxy && go run main.go
