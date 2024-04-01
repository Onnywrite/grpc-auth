GEN=./gen
PROTOPATH=./proto
SSOPROTOS=sso.proto sso_messages.proto

protoc:
	protoc --go_out=${GEN} --go_opt=paths=source_relative \
    --go-grpc_out=${GEN} --go-grpc_opt=paths=source_relative \
	--proto_path=${PROTOPATH} ${SSOPROTOS}

build:
	rm -rf bin
	go build -o ./bin/sso ./cmd/sso/main.go
	cp -r configs/ ./bin/
	cp -r storage/ ./bin/

run: build
	./bin/sso --config-path=./configs/local.yaml

pgadmin:
	sudo docker run --rm -dp 5000:80 \
	-e PGADMIN_DEFAULT_EMAIL=admin@gmail.com -e PGADMIN_DEFAULT_PASSWORD=admin \
	--network sso_db_test_pgnet_test dpage/pgadmin4

cover:
	go test -short -count=1 -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out

testdb_up:
	sudo docker-compose -f ./test-compose.yaml up -d

testdb_down:
	sudo docker-compose -f ./test-compose.yaml down