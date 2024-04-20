-include .env
export

CURRENT_DIR=$(shell pwd)
APP=dennic_user_service
CMD_DIR=./cmd

.DEFAULT_GOAL = build
POSTGRES_USER = postgres
POSTGRES_PASSWORD = 1234
POSTGRES_HOST = localhost
POSTGRES_PORT = 5544
POSTGRES_DATABASE = dennic_user_service

# build for current os
.PHONY: build
build:
	go build -ldflags="-s -w" -o ./bin/${APP} ${CMD_DIR}/app/main.go

# build for linux amd64
.PHONY: build-linux
build-linux:
	CGO_ENABLED=0 GOARCH="amd64" GOOS=linux go build -ldflags="-s -w" -o ./bin/${APP} ${CMD_DIR}/app/main.go

# run service
.PHONY: run
run:
	go run ${CMD_DIR}/app/main.go

# migrate
.PHONY: migrate
migrate:
	migrate -source file://migrations -database postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DATABASE}?sslmode=disable up

migrate-file:
	migrate create -ext sql -dir migrations/ -seq create_table_users

# proto
.PHONY: proto-gen
proto-gen:
	./scripts/gen-proto.sh

# git submodule init 	
.PHONY: pull-proto
pull-proto:
	git submodule update --init --recursive

# go generate	
.PHONY: go-gen
go-gen:
	go generate ./...

# run test
test:
	go test -v -cover -race ./internal/...

# -------------- for deploy --------------
build-image:
	docker build --rm -t ${REGISTRY}/${PROJECT_NAME}/${APP}:${TAG} .
	docker tag ${REGISTRY}/${PROJECT_NAME}/${APP}:${TAG} ${REGISTRY}/${PROJECT_NAME}/${APP}:${ENV_TAG}

push-image:
	docker push ${REGISTRY}/${PROJECT_NAME}/${APP}:${TAG}
	docker push ${REGISTRY}/${PROJECT_NAME}/${APP}:${ENV_TAG}

