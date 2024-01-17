.PHONY: build images

# Development

build_rev := "main"
ifneq ($(wildcard .git),)
	build_rev := $(shell git rev-parse --short HEAD)
endif
build_date := $(shell date -u '+%Y-%m-%dT%H:%M:%S')

setup:
	@go mod download

lint:
	@golangci-lint run --print-issued-lines=false --out-format=colored-line-number ./...

vet:
	@go vet ./...

test:
	@go test ./... -v


build:
	@go build -ldflags "-X main.commit=$(build_rev) -X main.date=$(build_date)" -o build/codapi -v cmd/main.go

run:
	@./build/codapi


# Containers

images:
	docker build --file images/alpine/Dockerfile --tag codapi/alpine:latest images/alpine/

network:
	docker network create --internal codapi

# Host OS

mount-tmp:
	mount -t tmpfs tmpfs /tmp -o rw,exec,nosuid,nodev,size=64m,mode=1777

# Deployment

app-download:
	@curl -L -o codapi.zip "https://api.github.com/repos/nalgeon/codapi/actions/artifacts/$(id)/zip"
	@unzip -ou codapi.zip
	@chmod +x build/codapi
	@rm -f codapi.zip
	@echo "OK"

app-start:
	@nohup build/codapi > codapi.log 2>&1 & echo $$! > codapi.pid
	@echo "started codapi"

app-stop:
	@kill $(shell cat codapi.pid)
	@rm -f codapi.pid
	@echo "stopped codapi"
