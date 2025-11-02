SERVICES := auth-service business-service consumer-service admin-service

.PHONY: b
b: 
	go build -o bin/auth-service ./backend/cmd/auth-service/main.go

.PHONY: u
u:
	docker-compose up --build

.PHONY: d
d:
	docker-compose down -v

.PHONY: a
a:
	docker exec -it ridehail-postgresql psql -U ridehail_user -W ridehail_db
