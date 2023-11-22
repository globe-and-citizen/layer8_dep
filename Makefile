## NPM Install
npm_install_all:
	cd sp_mock/frontend && npm install && cd ../backend && npm install  && cd ../../resource_server/frontend && npm install

go_mod_tidy_all:
	cd interceptor && go mod tidy && cd ../middleware && go mod tidy && cd ../proxy && go mod tidy && cd ../resource_server/backend && go mod tidy && cd ../../utils && go mod tidy

## Interceptor Calls
build_interceptor: ## must do from a bash terminal ..
	## Put WASM file directly in the CDN of the auth server
	## cd interceptor/ && GOARCH=wasm GOOS=js go build -o dist/interceptor.wasm && cp ./dist/interceptor.wasm ../proxy/assets/cdn/interceptor/interceptor__local.wasm && cp ./dist/wasm_exec.js ../proxy/assets/cdn/interceptor/wasm_exec.js
	## Put WASM file directly in the sp_mock frontend
	## cd interceptor/ && GOARCH=wasm GOOS=js go build -o dist/interceptor.wasm && cp ./dist/interceptor.wasm ../sp_mock/frontend/public/interceptor.wasm && cp ./dist/wasm_exec.js ../sp_mock/frontend/public/wasm_exec.js
	@'$(MAKE)' -C ./interceptor build

## Build Middleware
build_middleware:
	cd ./middleware/ && GOARCH=wasm GOOS=js go build -o ./dist/middleware.wasm

## Run Mock
run_frontend: # Port 5173
	cd sp_mock/frontend && npm run dev
	
run_backend: # Port 8000
	cd sp_mock/backend && npm run dev

## Run Proxy
run_proxy: # Port 5000
	cd proxy && go run main.go --server=proxy --port=5001

# Serve auth server
run_auth: # Port 5001
	cd proxy && go run main.go --server=auth

# Run Resource Server Backend
run_rs_backend: # Port 3000
	cd resource_server/backend && go run main.go

# Run Resource Server Frontend
run_rs_frontend: # Port 5174
	cd resource_server/frontend && npm run dev

run_db:
	docker run -d --rm \
		--name layer8-resource \
		-v $(PWD)/.docker/postgres:/var/lib/postgresql/data \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_PASSWORD=postgres \
		-e POSTGRES_DBNAME=postgres \
		-p 5434:5432 postgres:14.3
