# Avito sources config

Multi-account polling is configured through a JSON file on the server.

Recommended server path:

```sh
/var/www/avito/avito_sources.json
```

Enable it in `.env`:

```sh
AVITO_SOURCES_CONFIG=/var/www/avito/avito_sources.json
```

Example:

```json
{
  "sources": [
    {
      "enabled": true,
      "name": "Omoda Samara",
      "source": "polling: Проверка транспорта",
      "account_id": "375283938",
      "chat_id": "a2u-413531671-375283938",
      "client_id": "your_avito_client_id",
      "client_secret": "your_avito_client_secret"
    }
  ]
}
```

To add a new account/chat:

1. Add one object to `sources`.
2. Set `enabled` to `true`.
3. Fill `name`, `account_id`, `chat_id`, `client_id`, and `client_secret`.
4. Restart the service:

```sh
cd /var/www/avito
systemctl restart avito-monitor
```

Check that the source started:

```sh
journalctl -u avito-monitor -n 80 --no-pager
tail -80 logs/$(date +%F)-important.log
```

Each Matrix notification includes the source name, `account_id`, and `chat_id`.
