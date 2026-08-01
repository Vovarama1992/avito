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

```sh
cd /var/www/avito
./avito-monitor sources add \
  -name "Jaecoo Samara" \
  -account-id "ACCOUNT_ID" \
  -chat-id "CHAT_ID" \
  -client-id "CLIENT_ID" \
  -client-secret "CLIENT_SECRET"
systemctl restart avito-monitor
```

To show configured sources:

```sh
cd /var/www/avito
./avito-monitor sources list
```

To temporarily disable a source:

```sh
cd /var/www/avito
./avito-monitor sources disable -name "Jaecoo Samara"
systemctl restart avito-monitor
```

To enable it again:

```sh
cd /var/www/avito
./avito-monitor sources enable -name "Jaecoo Samara"
systemctl restart avito-monitor
```

Check that the source started:

```sh
journalctl -u avito-monitor -n 80 --no-pager
tail -80 logs/$(date +%F)-important.log
```

Each Matrix notification includes the source name, `account_id`, and `chat_id`.
