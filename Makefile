SERVICE_NAME=auth-service
ENGINE=main.go
BUILD_DIR=build

.PHONY: build run docker-build docker-run docker-build-run docker docker-down check-if-present-env check-if-valid-env clean-proto proto-gen ssl-gen

CHECK_ENV := production|staging|local

proto-gen:	clean-proto
	@echo "Generating the stubs"
	./scripts/proto-gen.sh
	@echo "Success generate stubs. All stubs created are in the 'stubs/' directory"
	@echo "Generating the Swagger UI"
	./scripts/swagger-ui-gen.sh
	@echo "Success generate Swagger UI. If you want to change Swagger UI to previous version copy the previous version from './cache/swagger-ui' directory"
	@echo "You can try swagger-ui with command 'make debug'"
	@echo "DO NOT EDIT ANY FILES STUBS!"

clean-proto:
	@echo "Delete all previous stubs ..."
	rm -rf stubs/*
	@echo "All stubs successfully deleted"

ssl-gen:
	@echo "Generating ssl configuration"
	./scripts/ssl-gen.sh
	@echo "Success generate ssl configuration. All SSL Configuration created in the 'ssl/' directory"
	@echo "DO NOT EXPOSE SSL DIRECTORY!"

docker-build: check-if-present-env check-if-valid-env
	@docker build . --file Dockerfile --build-arg ENVIRONMENT=${ENV} --no-cache --tag auth-service

docker-run:
	@docker run --name=auth-service -p 3001:3001 -d auth-service:latest
	@docker ps

docker-build-run: docker-build docker-run

docker:
	@echo Starting docker compose
	docker compose -f docker-compose.yaml up -d --build

docker-down:
	@echo Stopping docker compose
	docker compose -f docker-compose.yaml down

build:
	@echo "Building app"
	go build -o ${BUILD_DIR}/${SERVICE_NAME} ${ENGINE}
	@echo "Success build app. Your app is ready to use in 'build/' directory."

run:
	@go fmt ./...
	@go run .

# Environment test/check
check-if-present-env:
	test $(ENV)
	
check-if-valid-env:
ifneq ($(ENV), $(filter $(ENV), production staging local))
	$(error "ENV=$(CHECK_ENV)" is required)
endif