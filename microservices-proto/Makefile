PROTO_DIR := proto
OUT_DIR := golang

.PHONY: proto
proto:
	protoc -I $(PROTO_DIR) \
		--go_out $(OUT_DIR) --go_opt paths=source_relative \
		--go-grpc_out $(OUT_DIR) --go-grpc_opt paths=source_relative \
		$(shell find $(PROTO_DIR) -name '*.proto')

.PHONY: clean
clean:
	rm -rf $(OUT_DIR)
