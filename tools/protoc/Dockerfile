FROM --platform=linux/amd64 golang:latest

RUN apt-get update && \
    apt-get install -y unzip

# https://github.com/protocolbuffers/protobuf/releases
ENV PROTOC_VERSION=30.0
# https://github.com/protocolbuffers/protobuf-go/releases
ENV PROTOC_GEN_GO_VERSION=1.36.5
# https://github.com/grpc/grpc-go/tags
ENV PROTOC_GEN_GO_GPPC_VERSION=1.71.0

RUN curl -sfL "https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-linux-x86_64.zip" -o /tmp/protoc.zip

RUN curl -sfL "https://github.com/protocolbuffers/protobuf-go/releases/download/v${PROTOC_GEN_GO_VERSION}/protoc-gen-go.v${PROTOC_GEN_GO_VERSION}.linux.amd64.tar.gz" -o /tmp/protoc_gen_go.tar.gz && \
    git clone https://github.com/grpc/grpc-go && \
    cd grpc-go/cmd/protoc-gen-go-grpc && \
    git checkout tags/v${PROTOC_GEN_GO_GPPC_VERSION} && \
    go install && \
    cd - && \
    rm -rf grpc-go

RUN	cd /tmp && \
	tar -xzf /tmp/protoc_gen_go.tar.gz && \
    mv protoc-gen-go /usr/local/bin/protoc-gen-go && \
	unzip protoc.zip -d protoc && \
	ls -lah /tmp && \
	mv protoc/bin/protoc /usr/local/bin/ && \
	mv protoc/include/* /usr/local/include/ && \
	chmod +x /usr/local/bin/* && \
	rm -rf /tmp/* && \
    mkdir /workspace

WORKDIR /workspace

ENTRYPOINT ["make", "docker-generate-protoc"]