GATEWAY_SERVICE_BINARY=binary_file/gatewayServiceApp

run_postgres:
	@echo "stopping docker images (if running...)"
	docker-compose down
	@echo "building (when required) and starting docker images..."
	docker compose up postgres
	@echo "docker images built and started"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_gateway_service
	@echo "stopping docker images (if running...)"
	docker-compose down
	@echo "building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "docker images built and started"

## down: stop docker compose
down:
	@echo "stopping docker compose..."
	docker-compose down
	@echo "done"

## build_gateway_service: builds the tech task binary as a linux executable
build_gateway_service:
	@echo "building gateway service binary..."
	cd ../gateway-service && env GOOS=linux CGO_ENABLED=0 go build -o ${GATEWAY_SERVICE_BINARY} ./cmd/service
	@echo "done"
