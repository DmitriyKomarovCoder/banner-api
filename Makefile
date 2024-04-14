include .env
export

db:
	docker exec -it hammy-db psql -U $(DB_USER) -d $(DB_NAME)

doc:
	swag init -g cmd/app/main.go

cover:
	sh scripts/coverage_test.sh

clean_docker:
	docker-compose down
	docker system prune -af 
	docker volume prune -af
	docker system df

run:
	touch server.log
	docker-compose up -d