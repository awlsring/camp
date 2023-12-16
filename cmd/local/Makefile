SRC_DIR := $(shell git rev-parse --show-toplevel)
PKG_DIR := $(SRC_DIR)/pkg
REST_GEN_DIR := $(PKG_DIR)/gen/local
GRPC_GEN_DIR := $(PKG_DIR)/gen/local_grpc
APP_DIR := $(SRC_DIR)/cmd/local
BUILD_DIR := $(SRC_DIR)/build
CUR_DIR := $(shell pwd)
MODEL_DIR := $(SRC_DIR)/model/local
PROTO_DIR := $(MODEL_DIR)/build/proto
GRPC_GO_MOD := github.com/awlsring/local
OPEN_API_MODEL := $(SRC_DIR)/model/local/build/smithy/openapi/openapi/CampLocal.openapi.json

all: build

codegen:
	@echo "Building the local model"
	cd ${SRC_DIR}/model/local && smithy format model && smithy build
	@echo "Generating the local api go server and client"
	mkdir -p $(REST_GEN_DIR)
	go mod download golang.org/x/mod # tpm hack, fix this later
	go mod download golang.org/x/tools # tpm hack, fix this later
	go run github.com/ogen-go/ogen/cmd/ogen --target $(REST_GEN_DIR) -package camplocal --clean $(OPEN_API_MODEL)
		@echo "Generating the campd grpc server and go client"
	mkdir -p $(SRC_DIR)/gentmp
	mkdir -p $(PROTO_DIR)
	~/.local/share/coursier/bin/smithytranslate smithy-to-proto -i $(MODEL_DIR)/build/smithy/source/model $(PROTO_DIR)
	cd $(PROTO_DIR) && buf generate --template $(MODEL_DIR)/buf.gen.yaml --output $(SRC_DIR)/gentmp
	mkdir -p $(GRPC_GEN_DIR)
	mv $(SRC_DIR)/gentmp/local_grpc/awlsring/camp/local/* $(GRPC_GEN_DIR)
	mv $(SRC_DIR)/gentmp//local_grpc/awlsring/camp/common/* $(GRPC_GEN_DIR)
	mv $(SRC_DIR)/gentmp/local_grpc/aws/protocols/* $(GRPC_GEN_DIR)
	mv $(SRC_DIR)/gentmp/local_grpc/smithytranslate/* $(GRPC_GEN_DIR)
	mv $(SRC_DIR)/gentmp/local_grpc/smithy/framework/* $(GRPC_GEN_DIR)
	rm -rf $(SRC_DIR)/gentmp

build: codegen
	@echo "Building Camp Local app"
	go mod tidy
	@echo "Building Camp Local Rest API"
	go build -o $(BUILD_DIR)/local $(APP_DIR)/api/main.go

.PHONY: all