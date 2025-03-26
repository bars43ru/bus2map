GO_INTERNAL_PKG = "github.com/bars43ru/bus2map"
GO_FILES = $(shell find . -type f -name '*.go' -not -path "./api/*" | tr "\n" " ")
GO_FILES_GOLINES = $(shell find . -type f -name '*.go' -not -path "./api/*" | tr "\n" " ")

docker-generate-protoc: ### Internal command that is called from the docker container to generate pb
	rm -r ./api/bustracking/*** 2> /dev/null || (echo "dir ./api/bustracking/ was empty"; exit 0)
	protoc \
		--proto_path=. \
		--go_out=. \
		--go_opt=module='github.com/bars43ru/bus2map' \
		--go-grpc_out=. \
		--go-grpc_opt=module='github.com/bars43ru/bus2map' \
		./api/proto/*.proto

build-protoc:
	docker image rm 'bus2map-protoc:latest' 2> /dev/null || (echo "Image 'bus2map-protoc:latest' didn't exist so not removed."; exit 0)
	docker build -t go-gps2yandex-protoc ./tools/protoc

generate-protoc:
	docker run --rm \
			--mount src="$(shell pwd)",target=/workspace,type=bind \
			--workdir /workspace \
			go-gps2yandex-protoc:latest

deps-up:
	docker-compose -f ./build/local/docker-compose.yml up -d

deps-down:
	docker-compose -f ./build/local/docker-compose.yml down -v

lint: ### Run linter
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v1.64.8 golangci-lint run ./... -v

install-tools: ### Install the tools necessary to work on the project
	go install github.com/segmentio/golines@v0.12.2
	go install mvdan.cc/gofumpt@v0.7.0
	go install github.com/daixiang0/gci@v0.13.6

fmt: ### Code formatting
	golines --base-formatter=gofmt --ignore-generated -m 130 -w $(GO_FILES_GOLINES)
	gofumpt -w $(GO_FILES)
	gci write \
		--section Standard \
		--section Default \
		--section "Prefix($(GO_INTERNAL_PKG))" \
		$(GO_FILES) > /dev/null 2>&1