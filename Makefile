SERVICES := auth-service business-service consumer-service admin-service

.PHONY: b
b: 
	cd backend &&	go build -o ../bin/auth-service ./cmd/auth-service/main.go

.PHONY: u
u:
	docker compose up --build

.PHONY: d
d:
	docker compose down -v

.PHONY: a
a:
	docker exec -it bargo-postgres psql -U user -W bargo_db
