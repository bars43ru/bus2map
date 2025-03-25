docker-generate-protoc:
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

lint:
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v1.64.8 golangci-lint run ./... -v
