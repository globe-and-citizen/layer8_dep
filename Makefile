## NPM Install
npm_install_all:
	cd sp_mock/frontend && npm install && npm i layer8_interceptor && cd ../backend && npm install && npm i layer8_middleware

go_mod_tidy:
	cd ./server && go mod tidy

go_mod_tidy_all:
	cd interceptor && go mod tidy && cd ../middleware && go mod tidy && cd ../server && go mod tidy

go_test:
	cd server && go test ./... -v -cover

## Run Mock
run_frontend: # Port 5173
	cd sp_mock/frontend && npm run dev
	
run_backend: # Port 8000
	cd sp_mock/backend && npm run dev

# Serve 3-in-1 server
run_server: # Port 5001
	cd server && go run main.go

run_server_local: # Port 5001 with in-memory db
	cd server && go run main.go -port=5001 -jwtKey=secret -MpKey=secret -UpKey=secret

build_server_image:
	docker build --tag layer8-server --file Dockerfile .

build_sp_mock_frontend_image:
	cd sp_mock/frontend && docker build --tag sp_mock_frontend --file Dockerfile .

build_sp_mock_backend_image:
	cd sp_mock/backend && docker build --tag sp_mock_backend --file Dockerfile .

# To build all images at once
build_images:
	make build_server_image && make build_sp_mock_frontend_image && make build_sp_mock_backend_image

run_layer8_server_image:
	docker run -p 5001:5001 -t layer8-server

run_sp_mock_frontend_image:
	docker run -p 8080:8080 -t sp_mock_frontend

run_sp_mock_backend_image:
	docker run -p 8000:8000 -t sp_mock_backend

push_layer8_server_image:
	aws lightsail push-container-image --region ca-central-1 --service-name aws-container-service-t1 --label layer8-server-version-18 --image layer8-server:latest

push_sp_mock_frontend_image:
	aws lightsail push-container-image --region ca-central-1 --service-name container-service-2 --label frontendversion10 --image sp_mock_frontend:latest

push_sp_mock_backend_image:
	aws lightsail push-container-image --region ca-central-1 --service-name container-service-3 --label backendversion15 --image sp_mock_backend:latest

push_images:
	make push_layer8_server_image && make push_sp_mock_frontend_image && make push_sp_mock_backend_image

run_local_db:
	docker run -d --rm \
		--name layer8-resource \
		-v $(PWD)/.docker/postgres:/var/lib/postgresql/data \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_PASSWORD=postgres \
		-e POSTGRES_DBNAME=postgres \
		-p 5434:5432 postgres:14.3



## DEPRECATED CALLS
build_middleware:
	cd ./middleware/ && GOARCH=wasm GOOS=js go build -o ./dist/middleware.wasm && cp ./dist/middleware.wasm ../sp_mock/backend/dist/middleware.wasm


build_interceptor: ## must do from a bash terminal ..
	## Put WASM file directly in the CDN of the auth server
	## cd interceptor/ && GOARCH=wasm GOOS=js go build -o dist/interceptor.wasm && cp ./dist/interceptor.wasm ../proxy/assets/cdn/interceptor/interceptor__local.wasm && cp ./dist/wasm_exec.js ../proxy/assets/cdn/interceptor/wasm_exec.js
	## Put WASM file directly in the sp_mock frontend
	## cd interceptor/ && GOARCH=wasm GOOS=js go build -o dist/interceptor.wasm && cp ./dist/interceptor.wasm ../sp_mock/frontend/public/interceptor.wasm && cp ./dist/wasm_exec.js ../sp_mock/frontend/public/wasm_exec.js
	@'$(MAKE)' -C ./interceptor build

copy_wasm_exec_js:
	cp "$(shell go env GOROOT)/misc/wasm/wasm_exec.js" ./interceptor/dist/wasm_exec.js
	cp "$(shell go env GOROOT)/misc/wasm/wasm_exec.js" ./middleware/dist/wasm_exec.js
	cp "$(shell go env GOROOT)/misc/wasm/wasm_exec.js" ./server/assets-v1/cdn/interceptor/wasm_exec.js
	cp "$(shell go env GOROOT)/misc/wasm/wasm_exec.js" ./server/assets-v1/cdn/wasm_exec_v1.js

