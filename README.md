Migration

migrate -database "postgresql://root:secret@localhost:5434/rupay?sslmode=disable" -path db/migration up
migrate -database "postgresql://root:secret@localhost:5434/rupay?sslmode=disable" -path db/migration down

migrate create -ext sql -dir db/migration -seq add_price_to_article_transaction

SQLC

sqlc generate

SWAGGO

- regenerate the documentation
  swag init --parseDependency

test
