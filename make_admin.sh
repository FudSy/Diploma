#!/usr/bin/env bash

set -euo pipefail

LOGIN="${1:-}"
if [ -z "$LOGIN" ]; then
  echo "Ошибка: укажите логин пользователя."
  echo "Использование: $0 <login> [--docker]"
  exit 1
fi

MODE="${2:-}"

DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5433}"
DB_NAME="${DB_NAME:-test}"
DB_USER="${DB_USER:-postgres}"
DB_PASSWORD="${DB_PASSWORD:-postgres}"

SQL="UPDATE users SET role = 'ADMIN' WHERE login = '${LOGIN}'; SELECT login, role FROM users WHERE login = '${LOGIN}';"

if [ "$MODE" = "--docker" ]; then
  CONTAINER="${DOCKER_CONTAINER:-diploma-postgres}"
  echo "Подключение к контейнеру $CONTAINER..."
  RESULT=$(docker exec -i "$CONTAINER" \
    psql -U "$DB_USER" -d "$DB_NAME" -c \
    "UPDATE users SET role = 'ADMIN' WHERE login = '${LOGIN}'; SELECT login, role FROM users WHERE login = '${LOGIN}';")
else
  echo "Подключение к локальной БД ${DB_HOST}:${DB_PORT}/${DB_NAME}..."
  RESULT=$(PGPASSWORD="$DB_PASSWORD" psql \
    -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" \
    -c "UPDATE users SET role = 'ADMIN' WHERE login = '${LOGIN}';" \
    -c "SELECT login, role FROM users WHERE login = '${LOGIN}';")
fi

echo "$RESULT"

if echo "$RESULT" | grep -q "0 rows"; then
  echo ""
  echo "Пользователь с логином '${LOGIN}' не найден."
  exit 1
fi

echo ""
echo "Готово! Пользователь '${LOGIN}' теперь администратор."
