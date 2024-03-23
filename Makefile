repo_name := go-web-template

go-mod-init:
	go mod init $(repo_name)

.PHONY: wire
wire: ## Wire generate
	wire gen $(repo_name)/internal/setup/...

.PHONY: build
build: wire ## Build
	mkdir -p bin
	go build -mod=readonly -v -o bin/$(repo_name) ./cmd

.PHONY run:
run: build #
	#go run cmd/main.go
	./bin/$(repo_name)

.PHONY docker-build:
docker-build:
	docker build -t $(repo_name) .

.PHONY docker-run:
docker-run: docker-build
	docker run -p 8080:8080 $(repo_name)