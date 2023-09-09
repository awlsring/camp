SRC_DIR := $(shell git rev-parse --show-toplevel)
GEN_DIR := $(SRC_DIR)/generated
CUR_DIR := $(shell pwd)
GRADLE := $(SRC_DIR)/gradlew
LOCAL_SDK_DST := $(GEN_DIR)/camp_local
LOCAL_OPEN_API_MODEL := $(SRC_DIR)/model/local/build/smithyprojections/local/open-api/openapi/CampLocal.openapi.json

all: build

codegen:
	cd ${SRC_DIR}/model/local && gradle build && gradle build
	mkdir -p $(LOCAL_SDK_DST)
	go run github.com/ogen-go/ogen/cmd/ogen --target $(LOCAL_SDK_DST) -package camplocal --debug.noerr --clean $(LOCAL_OPEN_API_MODEL)

build: codegen

.PHONY: all