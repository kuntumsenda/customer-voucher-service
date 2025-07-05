.PHONY: proto clean

PROTO_DIR=proto
OUT_DIR=protogen

PROTOS=$(shell find $(PROTO_DIR) -name '*.proto')

proto:
	@echo "Compiling protos: $(PROTOS)"
	protoc -I=$(PROTO_DIR) \
		--go_out=$(OUT_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative \
		$(PROTOS)

clean:
	rm -rf $(OUT_DIR)/*