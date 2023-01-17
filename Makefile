get-docs:
	go get -u github.com/swaggo/swag/cmd/swag

docs: get-docs
	swag init --dir cmd/ --parseDependency --output docs

build:
	go build -o bin/restapi cmd/main.go

run:
	go run cmd/main.go
