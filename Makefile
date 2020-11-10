GORUN := go run cmd/main.go --conf=cmd/config.toml
MIGRATE := go run cmd/migrate/main.go --conf=cmd/config.toml
DOCKER := docker build -t gbbmn-cloud .
BUILD := cd cmd && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cloud
GENERATE := cd cmd/convert && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o pngpic

all:
	@echo 👌 start all

go:
	@echo 🚀 build
	@$(BUILD)

run:
	@echo 🚀 server
	@$(GORUN)

migrate:
	@echo 🚀 migrate
	@$(MIGRATE)

pic:
	@echo 🚀 generate picture
	@$(GENERATE)

docker:
	@echo 🚀 build in linux
	@$(BUILD)
	@echo 🚀 build in docker
	@$(DOCKER)

	
