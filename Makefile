IMAGE_NAME=bloxroute/bdn-protobuf:v3.19.3-serum

.PHONY: proto
proto: proto-build-api

.PHONY: proto-build-api
proto-build-api:
	docker run -v $(CURDIR)/proto:/go/protobuf/out \
			   -v $(CURDIR)/solana-trader-proto/proto:/go/protobuf/in $(IMAGE_NAME) \
		protoc --go_out=../out --go_opt=paths=source_relative  --go-grpc_out=../out --go-grpc_opt=paths=source_relative api.proto

.PHONY: test
test:
	go test -v ./...

.PHONY: unit
unit:
	 go test -v ./bxserum/provider_test/.
