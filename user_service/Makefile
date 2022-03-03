db_reset:
	docker-compose exec db psql -U postgres -c "DROP DATABASE IF EXISTS user_service"
	docker-compose exec db psql -U postgres -c "CREATE DATABASE user_service"
	go run cmd/initdb/main.go
