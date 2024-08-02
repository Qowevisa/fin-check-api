def: http-server
	@

http-server:
	go build -v -o ./bin/$@ ./cmd/$@

api-local: swagger http-server
	./bin/http-server

swagger:
	swag init -g ./cmd/http-server/main.go
