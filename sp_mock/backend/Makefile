include .env
export

build:
	@docker build -t layer8-mock-backend .

push:
	@docker tag layer8-mock-backend:latest registry.digitalocean.com/$(DIGITALOCEAN_REGISTRY)/layer8-mock-backend:latest
	@docker push registry.digitalocean.com/$(DIGITALOCEAN_REGISTRY)/layer8-mock-backend:latest
