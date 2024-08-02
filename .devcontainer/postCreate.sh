go version
go install github.com/tursodatabase/turso-cli/cmd/turso@latest
go install github.com/air-verse/air@latest
go install -tags 'postgres sqlite sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
