#!/bin/bash

Отправляет 12 быстрых запросов на /halloween и выводит HTTP-статус.

Ожидайте увидеть 200 OK, а затем 429 Too Many Requests.

ENDPOINT="http://localhost:8080/halloween"
REQUESTS=12

echo "--- Запуск теста Rate Limiting: $REQUESTS запросов ---"

Отправляем запросы в цикле

for i in $(seq 1 $REQUESTS); do

# shellcheck disable=SC2215
-s: бесшумный режим, -o /dev/null: игнорировать вывод тела, -w "%{http_code}": вывести только код состояния

STATUS=$(curl -s -o /dev/null -w "%{http_code}" $ENDPOINT)

if [ "$STATUS" -eq 429 ]; then
echo "Запрос №$i: $STATUS (TOO MANY REQUESTS) - УСПЕХ! ✅"
elif [ "$STATUS" -eq 200 ]; then
echo "Запрос №$i: $STATUS (OK)"
else
echo "Запрос №$i: $STATUS (Неожиданный статус)"
fi
done

echo "--- Тест завершен ---"