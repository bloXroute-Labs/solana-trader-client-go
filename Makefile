IMAGE_NAME=bloxroute/bdn-protobuf:v3.19.3-rp-gateway
LOCAL_IMAGE_NAME=bx-proto-gen

.PHONY: proto
proto: proto-build-api proto-build-gw

.PHONY: proto-build-gw
proto-build-gw:
	docker run -v $(CURDIR)/proto:/go/protobuf $(IMAGE_NAME) \
		protoc --grpc-gateway_out . --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative api.proto

.PHONY: proto-build-api
proto-build-api:
	docker run -v $(CURDIR)/proto:/go/protobuf $(IMAGE_NAME) \
		protoc --go_out=. --go_opt=paths=source_relative  --go-grpc_out=. --go-grpc_opt=paths=source_relative api.proto