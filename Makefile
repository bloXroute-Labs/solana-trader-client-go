IMAGE_NAME=bloxroute/bdn-protobuf:v3.19.3-rp-gateway

.PHONY: proto
proto: proto-build-api

.PHONY: proto-build-api
proto-build-api:
	docker run -v $(CURDIR)/proto:/go/protobuf $(IMAGE_NAME) \
		protoc --go_out=. --go_opt=paths=source_relative  --go-grpc_out=. --go-grpc_opt=paths=source_relative api.proto

.PHONY: test
test:
	go test -v ./...
