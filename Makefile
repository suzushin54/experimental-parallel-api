.PHONY: proto
proto:
	@which protoc > /dev/null || (echo "protoc is not installed" && exit 1)
	@which protoc-gen-go > /dev/null || go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@which protoc-gen-go-grpc > /dev/null || go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@find proto -name "*.proto" -exec protoc --proto_path=proto --go_out=gen --go_opt=paths=source_relative --go-grpc_out=gen --go-grpc_opt=paths=source_relative {} \;

.PHONY: clean
clean:
	@find ./gen -type f -name "*.go" -delete

.PHONY: gen


.PHONY: start
start:
	go run ./cmd/server/main.go
