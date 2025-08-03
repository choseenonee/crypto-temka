migration/up:
	goose -dir ./deploy/migrations postgres "user=salut password=salut123 dbname=root port=5433 host=localhost sslmode=disable" up

migration/down:
	goose -dir ./deploy/migrations postgres "user=salut password=salut123 dbname=root port=5433 host=localhost sslmode=disable" down


docs/gen:
	cd ./internal/delivery && swag init -g server.go --pd /Users/z.gabdrakhmanov/GolandProjects/crypto-temka/internal/models && cd ../../