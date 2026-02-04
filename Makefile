.PHONY: refresh full-refresh build up down logs app-logs commit

APP=avito_monitor_app

refresh:
	git pull origin main
	docker-compose build app
	docker-compose stop app
	docker-compose up -d --no-deps app
	docker-compose logs -f app

full-refresh:
	git pull origin main
	docker-compose down
	docker-compose build --no-cache
	docker-compose up -d
	docker-compose logs -f app

build:
	docker-compose build

up:
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs -f

app-logs:
	docker logs --tail=100 -f $(APP)

commit:
	git add .
	git commit -m "$${m:-update}"
	git push origin main