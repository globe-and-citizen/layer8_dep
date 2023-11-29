include .env
export

docker-build:
	@echo "Building proxy service..."
	@docker build -t layer8-proxy -f ./docker/proxy/Dockerfile .
	@echo "Building auth service..."
	@docker build -t layer8-auth -f ./docker/auth/Dockerfile .

docker-run:
	@echo "Running proxy service..."
	@docker run -d --network host --name layer8-proxy layer8-proxy
	@echo "Started proxy service on port 5001\n\n-----------------------------\n"
	@echo "Running auth service..."
	@docker run -d --network host --name layer8-auth layer8-auth
	@echo "Started auth service on port 5000"

docker-stop:
	@echo "Stopping proxy service..."
	@docker stop layer8-proxy && docker rm layer8-proxy
	@echo "Stopping auth service..."
	@docker stop layer8-auth && docker rm layer8-auth

docker-logs-proxy:
	@docker logs -f layer8-proxy

docker-logs-auth:
	@docker logs -f layer8-auth

docker-push-proxy: docker-push-proxy-do docker-push-proxy-hub

docker-push-auth: docker-push-auth-do docker-push-auth-hub

docker-push-proxy-do:
	@echo "pushing proxy service..."
	@docker tag layer8-proxy:latest registry.digitalocean.com/$(DIGITALOCEAN_REGISTRY)/layer8-proxy:latest
	@docker push registry.digitalocean.com/$(DIGITALOCEAN_REGISTRY)/layer8-proxy:latest
	@echo "pushed proxy service to registry.digitalocean.com/$(DIGITALOCEAN_REGISTRY)/layer8-proxy:latest"

docker-push-proxy-hub:
	@echo "pushing proxy service..."
	@docker tag layer8-proxy:latest $(DOCKER_HUB_REGISTRY)/layer8-proxy
	@docker push $(DOCKER_HUB_REGISTRY)/layer8-proxy
	@echo "pushed proxy service to $(DOCKER_HUB_REGISTRY)/layer8-proxy"

docker-push-auth-do:
	@echo "pushing auth service..."
	@docker tag layer8-auth:latest registry.digitalocean.com/$(DIGITALOCEAN_REGISTRY)/layer8-auth:latest
	@docker push registry.digitalocean.com/$(DIGITALOCEAN_REGISTRY)/layer8-auth:latest
	@echo "pushed auth service to registry.digitalocean.com/$(DIGITALOCEAN_REGISTRY)/layer8-auth:latest"

docker-push-auth-hub:
	@echo "pushing auth service..."
	@docker tag layer8-auth:latest $(DOCKER_HUB_REGISTRY)/layer8-auth
	@docker push $(DOCKER_HUB_REGISTRY)/layer8-auth
	@echo "pushed auth service to $(DOCKER_HUB_REGISTRY)/layer8-auth"
