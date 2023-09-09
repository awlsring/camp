SRC_DIR := $(shell git rev-parse --show-toplevel)
GEN_DIR := $(SRC_DIR)/generated
CUR_DIR := $(shell pwd)
MODEL_DIR := $(SRC_DIR)/models
GRADLE := $(SRC_DIR)/gradlew
LOCAL_SDK_DST := $(GEN_DIR)/camp_local
LOCAL_OPEN_API_MODEL := $(MODEL_DIR)/local/build/smithy/openapi/openapi/CampLocal.openapi.json

all: build

codegen:
	cd ${MODEL_DIR}/common && smithy build
	cd ${MODEL_DIR}/local && smithy build
	mkdir -p $(LOCAL_SDK_DST)
	go run github.com/ogen-go/ogen/cmd/ogen --target $(LOCAL_SDK_DST) -package camplocal --debug.noerr --clean $(LOCAL_OPEN_API_MODEL)

build: codegen

.PHONY: all