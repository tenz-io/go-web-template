repo_name := go-web-template

go-mod-init:
	go mod init $(repo_name)

.PHONY run:
run:
	go run cmd/main.go

.PHONY: wire
wire: ## Wire generate
	wire gen $(repo_name)/internal/setup/...

.docker-build:
docker-build:
	docker build -t $(repo_name) .

.docker-run:
docker-run: docker-build
	docker run -p 8080:8080 $(repo_name)