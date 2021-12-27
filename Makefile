VERSION := $(shell cat VERSION)
PROJECT := nudge
SOURCE_FILES ?=./internal/... ./cmd/...
BUILDS_PATH := ./build
IMAGE_PATH := fmartingr/nudge:${VERSION}

export SOURCE_FILES
export BUILDS_PATH

.PHONY: all
all: help

# this is godly
# https://news.ycombinator.com/item?id=11939200
.PHONY: help
help:	### this screen. Keep it first target to be default
ifeq ($(UNAME), Linux)
	@grep -P '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
else
	@# this is not tested, but prepared in advance for you, Mac drivers
	@awk -F ':.*###' '$$0 ~ FS {printf "%15s%s\n", $$1 ":", $$2}' \
		$(MAKEFILE_LIST) | grep -v '@awk' | sort
endif

.PHONY: build-docker
build-docker:
	$(info Make: Building images)
	@docker build -t ${IMAGE_PATH} .

.PHONY: run-docker
run-docker: build-docker
	$(info Make: Run docker)
	@docker run --rm -t ${IMAGE_PATH}

.PHONY: clean
clean: ###  clean test cache, build files
	$(info: Make: Clean)
	@rm -rf ./build
	@go clean ${CLEAN_OPTIONS}

.PHONY: build
build: clean ### builds the project for the setup os/arch combinations
	$(info: Make: Build)
	@go build -tags netgo -a -v -ldflags "${LD_FLAGS}" -o ${BUILDS_PATH}/nudge ./cmd/*.go
	@chmod +x ${BUILDS_PATH}/*

.PHONY: quick-run
quick-run: ### Executes the project using golang
	go run cmd/*.go

.PHONY: run
run: ### Executes the project build locally
	@make build
	${BUILDS_PATH}/nudge
