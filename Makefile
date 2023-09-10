SRC_DIR := $(shell git rev-parse --show-toplevel)
GEN_DIR := $(SRC_DIR)/generated
APP_DIR := $(SRC_DIR)/apps
CUR_DIR := $(shell pwd)
MODEL_DIR := $(SRC_DIR)/models
GRADLE := $(SRC_DIR)/gradlew
LOCAL_OPEN_API_MODEL := $(MODEL_DIR)/local/build/smithy/openapi/openapi/CampLocal.openapi.json
LOCAL_SDK_DST := $(GEN_DIR)/camp_local

all: build

codegen: codegen-common codegen-local
	@echo "All codegen complete"

codegen-common:
	@echo "Building the common model"
	cd ${MODEL_DIR}/common && smithy build

codegen-local: codegen-common
	@echo "Building the local model"
	cd ${MODEL_DIR}/local && smithy build
	@echo "Generating the local api server and client"
	mkdir -p $(LOCAL_SDK_DST)
	go run github.com/ogen-go/ogen/cmd/ogen --target $(LOCAL_SDK_DST) -package camplocal --debug.noerr --clean $(LOCAL_OPEN_API_MODEL)
	@echo "Creating static docs"
	mkdir -p $(APP_DIR)/local/swagger
	cp $(LOCAL_OPEN_API_MODEL) $(APP_DIR)/local/swagger/swagger.json

build: codegen

build-local: codegen-local

.PHONY: all