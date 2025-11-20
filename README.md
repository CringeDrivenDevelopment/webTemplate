![Coverage](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/HeyReyHR/33fad6beea65140d890305e435b5b0ab/raw/coverage.json)

# Backend Template

```shell
# Сгенерировать Swagger документацию
task swag

# Запуск Unit тестов
task test

# Запустить линтер
task lint

# Сгенерировать код от SQLC
task sqlc

# Dev среда
task dcu
task dcd
```

## E2E тесты
```shell
# tests/e2e
uv run pytest -n auto
```