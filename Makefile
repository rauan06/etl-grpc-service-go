# Makefile for generating Go code from proto files

# Set default proto directory (can be overridden by make command)
PROTO_DIR ?= etl

# Specify the proto files and output directory based on PROTO_DIR
PROTO_FILES = protos/$(PROTO_DIR)/v1/*.proto
OUT_DIR = ./protos/$(PROTO_DIR)/v1/pb

# Go plugin options
GO_OPTS = --go_out=$(OUT_DIR) --go_opt=paths=source_relative
GO_GRPC_OPTS = --go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative

# Proto include paths
GOOGLEAPIS_DIR = googleapis
INCLUDE_PATHS = -I. -I$(GOOGLEAPIS_DIR)

# Target to generate Go code from proto files
.PHONY: generate
generate:
	@echo "Generating Go code for $(PROTO_DIR) proto files..."
	protoc $(INCLUDE_PATHS) $(PROTO_FILES) $(GO_OPTS) $(GO_GRPC_OPTS)

# Add a clean target to remove generated files (optional)
.PHONY: clean
clean:
	@echo "Cleaning up generated files in $(OUT_DIR)..."
	rm -rf $(OUT_DIR)
