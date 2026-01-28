#!/bin/bash

# Sends 12 fast requests to /halloween and outputs the HTTP status.
#
# Expect to see 200 OK, then 429 Too Many Requests.

ENDPOINT="http://localhost:8080/halloween"
REQUESTS=12

echo "--- Starting Rate Limiting Test: $REQUESTS requests ---"

# Sending requests in a loop
for i in $(seq 1 $REQUESTS); do

  # shellcheck disable=SC2215
  # -s: silent mode, -o /dev/null: ignore body output, -w "%{http_code}": output only status code
  STATUS=$(curl -s -o /dev/null -w "%{http_code}" $ENDPOINT)

  if [ "$STATUS" -eq 429 ]; then
    echo "Request #$i: $STATUS (TOO MANY REQUESTS) - SUCCESS! ✅"
  elif [ "$STATUS" -eq 200 ]; then
    echo "Request #$i: $STATUS (OK)"
  else
    echo "Request #$i: $STATUS (Unexpected status)"
  fi
done

echo "--- Test finished ---"
