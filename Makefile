## Interceptor Calls
build_interceptor: ## must do from a bash terminal
	cd interceptor/ && GOARCH=wasm GOOS=js go build -o dist/interceptor.wasm && cp ./dist/interceptor.wasm ../sp_mock/frontend/public/interceptor.wasm && cp ./dist/wasm_exec.js ../sp_mock/frontend/public/wasm_exec.js

## Build Middleware
build_middleware:
	cd ./middleware/ && GOARCH=wasm GOOS=js go build -o ./dist/middleware.wasm

## Run Mock
run_mock:
	make run_mock_frontend && run_mock_backend

run_mock_frontend:
	cd sp_mock/frontend && npm run dev
	
run_mock_backend:
	cd sp_mock/backend && npm run dev

## Run Proxy
build-layer8-slave-one:
	cd layer8/proxy_slave/layer8-slave-one/cmd && go build -o ../../../bin/layer8-slave-one

layer8-slave-one: # Port 8001
	make build-layer8-slave-one && ./layer8/bin/layer8-slave-one

generate-layer8-slave-proto:
	cd go-layer8-slaves && protoc --go_out=. --go-grpc_out=. proto/Layer8Slave.proto