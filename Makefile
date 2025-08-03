migration/local/up:
	goose -dir ./deploy/migrations postgres "user=salut password=salut123 dbname=root port=5433 host=localhost sslmode=disable" up

migration/local/down:
	goose -dir ./deploy/migrations postgres "user=salut password=salut123 dbname=root port=5433 host=localhost sslmode=disable" down

migration/remote/up:
	goose -dir ./deploy/migrations postgres "user=salut password=salut123 dbname=root port=5433 host=5.129.198.124 sslmode=disable" up

migration/remote/down:
	goose -dir ./deploy/migrations postgres "user=salut password=salut123 dbname=root port=5433 host=5.129.198.124 sslmode=disable" down


docs/gen:
	cd ./internal/delivery && swag init -g server.go --pd /Users/z.gabdrakhmanov/GolandProjects/crypto-temka/internal/models && cd ../../