SRC_DIR := $(shell git rev-parse --show-toplevel)
PKG_DIR := $(SRC_DIR)/pkg
REST_GEN_DIR := $(PKG_DIR)/gen/campd
GRPC_GEN_DIR := $(PKG_DIR)/gen/campd_grpc
APP_DIR := $(SRC_DIR)/cmd/campd
BUILD_DIR := $(SRC_DIR)/build
CUR_DIR := $(shell pwd)
MODEL_DIR := $(SRC_DIR)/model/campd
PROTO_DIR := $(MODEL_DIR)/build/proto
GRPC_GO_MOD := github.com/awlsring/campd
OPEN_API_MODEL := $(MODEL_DIR)/build/smithy/openapi/openapi/Campd.openapi.json

all: build

codegen:
	@echo "Building the campd model"
	cd ${SRC_DIR}/model/campd && smithy format model && smithy build
	@echo "Generating the campd rest server and go client"
	mkdir -p $(REST_GEN_DIR)
	go mod download golang.org/x/mod # tpm hack, fix this later
	go mod download golang.org/x/tools # tpm hack, fix this later
	go run github.com/ogen-go/ogen/cmd/ogen --target $(REST_GEN_DIR) -package campd --clean $(OPEN_API_MODEL)
	@echo "Generating the campd grpc server and go client"
	mkdir -p $(SRC_DIR)/gentmp
	mkdir -p $(PROTO_DIR)
	~/.local/share/coursier/bin/smithytranslate smithy-to-proto -i $(MODEL_DIR)/build/smithy/source/model $(PROTO_DIR)
	cd $(PROTO_DIR) && buf generate --template $(MODEL_DIR)/buf.gen.yaml --output $(SRC_DIR)/gentmp
	mkdir -p $(PKG_DIR)/gen/campd_grpc
	mv $(SRC_DIR)/gentmp/campd_grpc/awlsring/camp/campd/* $(PKG_DIR)/gen/campd_grpc
	mv $(SRC_DIR)/gentmp//campd_grpc/awlsring/camp/common/* $(PKG_DIR)/gen/campd_grpc
	mv $(SRC_DIR)/gentmp/campd_grpc/aws/protocols/* $(PKG_DIR)/gen/campd_grpc
	mv $(SRC_DIR)/gentmp/campd_grpc/smithytranslate/* $(PKG_DIR)/gen/campd_grpc
	rm -rf $(SRC_DIR)/gentmp

build: codegen
	@echo "Building Campd app"
	go mod tidy
	@echo "Building Camp Local Rest API"
	go build -o $(BUILD_DIR)/campd $(APP_DIR)/main.go

.PHONY: all