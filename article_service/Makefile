db_reset:
	docker-compose exec db psql -U postgres -c "DROP DATABASE IF EXISTS article_service"
	docker-compose exec db psql -U postgres -c "CREATE DATABASE article_service"
	go run cmd/initdb/main.go
