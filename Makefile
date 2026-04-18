DB_URL=postgres://postgres:postgres@localhost:5432/messenger?sslmode=disable

migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

migrate-up:
	migrate -path migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path migrations -database "$(DB_URL)" down 1

migrate-force:
	migrate -path migrations -database "$(DB_URL)" force $(version)

migrate-version:
	migrate -path migrations -database "$(DB_URL)" version

.PHONY: migrate-create migrate-up migrate-down migrate-force migrate-version