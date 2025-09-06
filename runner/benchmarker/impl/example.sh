#!/bin/sh

if [ -z "$1" ]; then
  echo "Usage: $0 <target URL>"
  exit 1
fi

TARGET_URL=$1

i=1
while [ $i -le 20 ]; do
  echo "Score: $(($i * 100))"
  if [ $(($i % 5)) -eq 0 ]; then
    echo "Error: Some error" >&2
  fi
  sleep 1
  i=$((i + 1))
done
