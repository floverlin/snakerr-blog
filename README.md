# Переменные окружения
```
SECRET=secret_string
```

# Флаги
```
-p PORT integer default=8080
-e ENVIROMENT "dev" | "prod" default="dev"
-c CONFIG_PATH string default="config.json"
```

# Переменные для goose
```
GOOSE_DRIVER=sqlite
GOOSE_DBSTRING=./database/data.db
GOOSE_MIGRATION_DIR=./migrations
```
