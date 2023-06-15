postgres:
	sudo docker run --name postgres -p 5435:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:alpine

createbd: 
	sudo docker exec -it postgres createdb --username=root --owner=root admins
	
dropdb:
	sudo docker exec -it postgres dropdb  students

migrateup: 
	migrate -path migration/ -database "postgresql://root:secret@localhost:5435/admins?sslmode=disable" -verbose up

migratedown:
	migrate -path migration/ -database "postgresql://root:secret@localhost:5435/admins?sslmode=disable" -verbose down

proto:
	protoc --go_out=./pkg/admin --go_opt=paths=source_relative \
    --go-grpc_out=./api/proto --go-grpc_opt=paths=source_relative \
    api/proto/service.proto



.PHONY: postgres createbd migrateup run
