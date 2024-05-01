.PHONY: run docker-build docker-run docker-build-run check-if-present-env check-if-valid-env

CHECK_ENV := production|staging|local

docker-build: check-if-present-env check-if-valid-env
	@docker build . --file Dockerfile --build-arg ENVIRONMENT=${ENV} --no-cache --tag auth-service

docker-run:
	@docker run --name=auth-service -p 3001:3001 -d auth-service:latest
	@docker ps

docker-build-run: docker-build docker-run

local:
	@cp .env.local .env
	@go build .
	@go run .

# Environment test/check
check-if-present-env:
	test $(ENV)
	
check-if-valid-env:
ifneq ($(ENV), $(filter $(ENV), production staging local))
	$(error "ENV=$(CHECK_ENV)" is required)
endif