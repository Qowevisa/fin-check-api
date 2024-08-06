def: http-server
	@

http-server:
	go build -v -o ./bin/$@ ./cmd/$@

api-local: swagger http-server
	./bin/http-server

swagger:
	#swag init -g ./cmd/http-server/main.go
	swag init -d ./cmd/http-server/,./types/,./handlers,./db

cloc-api:
	cloc ./cmd/http-server/ ./db/ ./handlers/ ./types/ ./middleware/ ./utils/ ./tokens/

install-swag:
	go install github.com/swaggo/swag/cmd/swag@latest
