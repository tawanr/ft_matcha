include .env
export

db-shell:
	sqlite3 $(DB_HOST)

migration:
	migrate create -ext sql -dir ./migrations -seq $(name)

migrate:
	migrate -path=./migrations -database=sqlite3://data.db up

init-db:
	go run internal/tools/initDB.go path=$(path)

test:
	go test ./...

tailwind:
	tailwindcss -i ui/static/css/base_input.css -o ui/static/css/base.css --watch
