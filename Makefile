format:
	gofmt -s -w .
	
test:
	go test -v -race ./...

dep-download:
	env GO111MODULE=on go mod download

tidy:
	env GO111MODULE=on go mod tidy

vendor:
	env GO111MODULE=on go mod vendor

cover:
	go test -v -race ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	go tool cover -func coverage.out 

coverhtml:
	go test -v -race ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

migration:
	# https://github.com/golang-migrate/migrate
	migrate create -ext sql -dir db/migrations -seq $(name)

migrate:
	# https://github.com/golang-migrate/migrate
	migrate -path db/migrations -database $(url) up

rollback:
	# https://github.com/golang-migrate/migrate
	migrate -path db/migrations -database $(url) down 1

force-migrate:
	# https://github.com/golang-migrate/migrate
	migrate -path db/migrations -database $(url) force $(version)