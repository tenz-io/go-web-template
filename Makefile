repo_name := go-web-template

TOOL_PATH = tool
TOOL_FILES := $(shell cd $(TOOL_PATH); find . -maxdepth 1 -type d|grep -v '/common')
TOOL_TAGET := $(basename $(patsubst ./%,%,$(TOOL_FILES)))
TOOL_BIN_FILES := $(basename $(patsubst ./%,$(bin_dir)/%,$(TOOL_FILES)))

go-mod-init:
	go mod init $(repo_name)

.PHONY: generate
generate: ## Generate
	go generate ./internal/constant/... ./internal/repository/... ./internal/service/...

.PHONY: wire
wire: ## Wire generate
	wire gen $(repo_name)/internal/setup/...

.PHONY: gci
gci:
	gci write -s standard -s default -s "prefix(github.com)" -s "prefix(go-web-template)" --skip-generated *


.PHONY: proto-compile
proto-compile:
	@echo "=== generate gin api"
	protoc -I ./api \
	--openapiv2_out ./api --openapiv2_opt logtostderr=true \
	--openapiv2_opt json_names_for_fields=false \
	--go_out ./api --go_opt=paths=source_relative \
	--go-gin_out ./api --go-gin_opt=paths=source_relative \
	./api/http/app/*.proto

	@echo "=== generate tag for gin api"
	protoc-go-inject-tag -input=./api/http/app/*.pb.go

.PHONY: build
build: wire ## Build
	mkdir -p bin
	go build -mod=readonly -v -o bin/$(repo_name) ./cmd

.PHONY: dep
dep:
	go mod tidy -v


.PHONY: test
test:
	go test -v ./... -cover

.PHONY run:
run: build #
	#go run cmd/main.go
	./bin/$(repo_name) -c config/app.yaml -p 8090 -v

$(TOOL_TAGET):
	@echo "=== build tool $@"
	scripts/build-tool.sh $@

.PHONY: build-tools
build-tools: $(TOOL_TAGET)

.PHONY docker-build:
docker-build:
	docker build --build-arg SKAFFOLD_GO_GCFLAGS="-N -l" -t $(repo_name) .

.PHONY docker-run:
docker-run: docker-build
	docker rm -f $(repo_name) || true
	docker run --name $(repo_name) -p 8080:8080 -p 8085:8085 $(repo_name)

