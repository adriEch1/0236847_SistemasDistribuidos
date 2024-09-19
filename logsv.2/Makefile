CONFIG_PATH=$(HOME)/.prolog/

.PHONY: init
init:
	mkdir -p $(CONFIG_PATH)

.PHONY: gencert
gencert:
	cfssl gencert \
		-initca test/ca-csr.json | cfssljson -bare ca

	cfssl gencert \
		-ca=ca.pem \
		-ca-key=ca-key.pem \
		-config=test/ca-config.json \
		-profile=server \
		test/server-csr.json | cfssljson -bare server
	mv *.pem *.csr $(CONFIG_PATH)

compile:
	protoc api/v1/*.proto \
		--go_out=./ \
		--go_opt=paths=source_relative \
		--proto_path=.

test:
	go test -race ./...

compile_rpc:
	protoc api/v1/*.proto \
		--go_out=./ \
		--go_opt=paths=source_relative \
		--proto_path=. \
		--go-grpc_out=./ \
		--go-grpc_opt=paths=source_relative
