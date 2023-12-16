SRC_DIR := $(shell git rev-parse --show-toplevel)
PKG_DIR := $(SRC_DIR)/pkg
GEN_DIR := $(PKG_DIR)/gen/local
APP_DIR := $(SRC_DIR)/cmd/local
BUILD_DIR := $(SRC_DIR)/build
CUR_DIR := $(shell pwd)
OPEN_API_MODEL := $(SRC_DIR)/model/local/build/smithy/openapi/openapi/CampLocal.openapi.json

all: build

codegen:
	@echo "Building the local model"
	cd ${SRC_DIR}/model/local && smithy format model && smithy build
	@echo "Generating the local api go server and client"
	mkdir -p $(GEN_DIR)
	go mod download golang.org/x/mod # tpm hack, fix this later
	go mod download golang.org/x/tools # tpm hack, fix this later
	go run github.com/ogen-go/ogen/cmd/ogen --target $(GEN_DIR) -package camplocal --clean $(OPEN_API_MODEL)

build: codegen
	@echo "Building Camp Local app"
	go mod tidy
	@echo "Building Camp Local Rest API"
	go build -o $(BUILD_DIR)/local $(APP_DIR)/api/main.go

.PHONY: all