include .env

run:
	go run ./cmd/api 

migrateup:
	migrate -path migrations -database "$(DB_TYPE)://$(DB_USER):$(DB_PASSWORD)@tcp($(DB_HOST))/$(DB_NAME)?parseTime=true" up

migratedown:
	migrate -path migrations -database "$(DB_TYPE)://$(DB_USER):$(DB_PASSWORD)@tcp($(DB_HOST))/$(DB_NAME)?parseTime=true" down