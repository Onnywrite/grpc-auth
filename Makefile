GEN=./gen
PROTOPATH=./proto
SSOPROTO=sso.proto

protoc:
	protoc --go_out=${GEN} --go_opt=paths=source_relative \
    --go-grpc_out=${GEN} --go-grpc_opt=paths=source_relative \
	--proto_path=${PROTOPATH} ${SSOPROTO}

build:
	rm -rf bin
	go build -o ./bin/sso ./cmd/sso/main.go
	ln -s ./configs configs
	ln -s ./storage storage

run: build
	./bin/sso