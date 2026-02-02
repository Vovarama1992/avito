.PHONY: refresh full-refresh build up down logs app-logs commit

APP=avito_app

# --- быстрый деплой ---
refresh:
	git pull origin master
	docker compose build app
	docker compose stop app
	docker compose up -d --no-deps app
	docker compose logs -f app

# --- полный рефреш ---
full-refresh:
	git pull origin master
	docker compose down
	docker compose build --no-cache
	docker compose up -d
	docker compose logs -f app

# --- сборка контейнеров ---
build:
	docker compose build

# --- поднять ---
up:
	docker compose up -d

# --- остановить ---
down:
	docker compose down

# --- логи всех сервисов ---
logs:
	docker compose logs -f

# --- логи приложения ---
app-logs:
	docker logs --tail=100 -f $(APP)

# --- Git ---
commit:
	git add .
	git commit -m "$${m:-update}"
	git push origin master