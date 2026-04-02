.PHONY: help generate-mocks test with-in-memo-docker-up with-in-memo-docker-down with-postgres-docker-up with-postgres-docker-down

help:
	@echo generate-mocks - Генерация всех моков в проекте
	@echo test - Запуск всех unit-тестов
	@echo with-in-memo-docker-up - Сборка приложения с реализацией in-memory storage
	@echo with-in-memo-docker-down - Остановка контейнеров in-memory storage
	@echo with-postgres-docker-up - Сборка приложения с реализацией PostgreSQL storage
	@echo with-postgres-docker-down - Остановка контейнеров PostgreSQL storage

generate-mocks:
	@echo Generating mocks...
	go generate ./...
	@echo Successfully generated

test:
	@echo Running all unit-tests...
	go test ./... -cover -v

with-in-memo-docker-up:
	@echo Building Docker Image with In-Memory...
	docker compose -f docker-compose.yml up --build -d

with-in-memo-docker-down:
	@echo Stopping In-Memory containers...
	docker compose -f docker-compose.yml down

with-postgres-docker-up:
	@echo Building Docker Image with PostgreSQL...
	docker compose -f docker-compose.postgres.yml up --build -d

with-postgres-docker-down:
	@echo Stopping PostgreSQL containers...
	docker compose -f docker-compose.postgres.yml down

