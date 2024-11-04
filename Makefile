.PHONY: proto
proto:
	@which protoc > /dev/null || (echo "protoc is not installed" && exit 1)
	@which protoc-gen-go > /dev/null || go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@which protoc-gen-go-grpc > /dev/null || go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	protoc --go_out=. --go-grpc_out=. proto/payment.proto

.PHONY: clean
clean:
	rm -f proto/*.go

.PHONY: start
start:
	go run ./cmd/server/main.go
