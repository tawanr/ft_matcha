include .env
export

db-shell:
	turso db shell http://db:8080

migration:
	migrate create -ext sql -dir ./migrations -seq $(name)

migrate:
	migrate -path=./migrations -database=sqlite3://data.db up

init-db:
	go run internal/tools/initDB.go path=$(path)

tailwind:
	tailwindcss -i ui/static/css/base_input.css -o ui/static/css/base.css --watch
