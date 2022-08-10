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

.PHONY: grpc-examples
grpc-examples:
	 go run ./examples/grpcclient/main.go

.PHONY: http-examples
http-examples:
	 go run ./examples/httpclient/main.go

.PHONY: ws-examples
ws-examples:
	 go run ./examples/wsclient/main.go

.PHONY: ssl-testnet ssl-mainnet cred-github environment-dev

environment-dev: ssl-testnet cred-github

ssl-testnet:
	mkdir -p $(CURDIR)/ssl/testnet/bloxroute_cloud_api/registration_only
	aws s3 cp s3://internal-credentials.bxrtest.com/bloxroute_cloud_api/registration_only/bloxroute_cloud_api_cert.pem $(CURDIR)/ssl/testnet/bloxroute_cloud_api/registration_only/
	aws s3 cp s3://internal-credentials.bxrtest.com/bloxroute_cloud_api/registration_only/bloxroute_cloud_api_key.pem $(CURDIR)/ssl/testnet/bloxroute_cloud_api/registration_only/

ssl-mainnet:
	mkdir -p $(CURDIR)/ssl/bloxroute_cloud_api/registration_only
	aws s3 cp s3://internal-credentials.blxrbdn.com/bloxroute_cloud_api/registration_only/bloxroute_cloud_api_cert.pem $(CURDIR)/ssl/bloxroute_cloud_api/registration_only/
	aws s3 cp s3://internal-credentials.blxrbdn.com/bloxroute_cloud_api/registration_only/bloxroute_cloud_api_key.pem $(CURDIR)/ssl/bloxroute_cloud_api/registration_only/

cred-github:
	aws s3 cp s3://files.bloxroute.com/trader-api/.netrc $(CURDIR)/.netrc