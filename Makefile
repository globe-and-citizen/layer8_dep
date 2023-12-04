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
run_proxy: # Port 5001
	cd proxy && go run main.go --server=proxy --port=5001
# cd proxy && go build -o bin/proxy && ./bin/proxy --server=proxy --port=5001

# Serve auth server
run_auth: # Port 5000
	cd proxy && go run main.go

# Run Resource Server Backend
run_rs_backend: # Port 3050
	cd resource_server/backend && go run main.go

# Run Resource Server Frontend
run_rs_frontend: # Port 5174
	cd resource_server/frontend && npm run dev

# To build all images at once
build_images:
	docker build --tag proxy-server --file Dockerfile_server .
	docker build --tag auth-server --file Dockerfile_auth .
	docker build --tag resource-server --file Dockerfile_resourceServer .

# To push all images to AWS Lightsail at once
push_images:
	aws lightsail push-container-image --region ca-central-1 --service-name aws-container-service-t1 --label proxy-server-v1 --image proxy-server:latest
	aws lightsail push-container-image --region ca-central-1 --service-name aws-container-service-t1 --label auth-server-v1 --image auth-server:latest
	aws lightsail push-container-image --region ca-central-1 --service-name aws-container-service-t1 --label resource-server-v1 --image resource-server:latest

build_proxy_server_image:
	docker build --tag proxy-server --file Dockerfile_server .

build_auth_server_image:
	docker build --tag auth-server --file Dockerfile_auth .

build_resource_server_image:
	docker build --tag resource-server --file Dockerfile_resourceServer .

run_proxy_server_image:
	docker run -p 5001:5001 proxy-server

run_auth_server_image:
	docker run -p 5000:5000 auth-server

run_resource_server_image:
	docker run -p 3050:3050 resource-server

push_proxy_server_image:
	aws lightsail push-container-image --region ca-central-1 --service-name aws-container-service-t1 --label proxy-server-v1 --image proxy-server:latest

push_auth_server_image:
	aws lightsail push-container-image --region ca-central-1 --service-name aws-container-service-t1 --label auth-server-v1 --image auth-server:latest

push_resource_server_image:
	aws lightsail push-container-image --region ca-central-1 --service-name aws-container-service-t1 --label resource-server-v1 --image resource-server:latest
