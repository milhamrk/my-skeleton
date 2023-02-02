APP_BIN = app/build/app

migrate:
	@./bin/migrate create -ext sql -dir migrations -seq -digits 8 $(NAME)

migrate.up:
	migrate -path migrations -database ${DATABASE_URL} up
