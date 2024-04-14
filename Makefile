include .env
export

db:
	docker exec -it hammy-db psql -U $(DB_USER) -d $(DB_NAME)

run:
	touch server.log
	docker-compose up -d