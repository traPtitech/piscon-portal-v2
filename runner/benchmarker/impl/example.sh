#!/bin/bash

if [ -z "$1" ]; then
  echo "Usage: $0 <target URL>"
  exit 1
fi

TARGET_URL=$1

for i in {1..20}; do 
  echo "Score: $(($i * 100))"
  if [ $(($i % 5)) -eq 0 ]; then
    echo "Error: Some error" >&2
  fi
  sleep 1
done
