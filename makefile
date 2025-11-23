# Переменные
DOCKER_COMPOSE = docker-compose -f docker/docker-compose.yml

.PHONY: all build up down logs restart clean

# Команда по умолчанию
all: up

# Сборка и запуск контейнеров
up:
	$(DOCKER_COMPOSE) up --build -d

# Остановка контейнеров
down:
	$(DOCKER_COMPOSE) down

# Просмотр логов
logs:
	$(DOCKER_COMPOSE) logs -f

# Рестарт
restart: down up

# Очистка (включая volume с базой данных)
clean:
	$(DOCKER_COMPOSE) down -v