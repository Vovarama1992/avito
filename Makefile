.PHONY: refresh full-refresh build up down logs app-logs status restart commit matrix-env

APP=avito-monitor
SERVICE=avito-monitor

refresh:
	git pull origin main
	/usr/local/go/bin/go mod download
	/usr/local/go/bin/go build -o $(APP) ./cmd
	systemctl restart $(SERVICE)
	journalctl -u $(SERVICE) -f

full-refresh:
	git pull origin main
	/usr/local/go/bin/go clean -cache
	/usr/local/go/bin/go mod download
	/usr/local/go/bin/go build -o $(APP) ./cmd
	systemctl restart $(SERVICE)
	journalctl -u $(SERVICE) -f

build:
	/usr/local/go/bin/go mod download
	/usr/local/go/bin/go build -o $(APP) ./cmd

up:
	systemctl enable $(SERVICE)
	systemctl start $(SERVICE)

down:
	systemctl stop $(SERVICE)

restart:
	systemctl restart $(SERVICE)

status:
	systemctl status $(SERVICE)

logs:
	journalctl -u $(SERVICE) -f

app-logs:
	journalctl -u $(SERVICE) -n 100 -f

matrix-env:
	nano .env

commit:
	git add .
	git commit -m "$${m:-update}"
	git push origin main