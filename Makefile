## NPM install
npm_install_all:
	cd sp_mock/frontend && npm install && cd ../backend && npm install 

## Interceptor Calls
build_interceptor: ## must do from a bash terminal
	cd interceptor/ && GOARCH=wasm GOOS=js go build -o dist/interceptor.wasm && cp ./dist/interceptor.wasm ../sp_mock/frontend/public/interceptor.wasm && cp ./dist/wasm_exec.js ../sp_mock/frontend/public/wasm_exec.js

## Build Middleware
build_middleware:
	cd ./middleware/ && GOARCH=wasm GOOS=js go build -o ./dist/middleware.wasm

## Run Mock

run_frontend:
	cd sp_mock/frontend && npm run dev
	
run_backend:
	cd sp_mock/backend && npm run dev

## Run Proxy

run_proxy:
	cd proxy && go run main.go


