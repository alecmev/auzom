#!/bin/bash

read -p "Rebuild migrations image? [yN] " build
if [ "$build" == "y" ]; then
  make migrations
fi

docker run --rm --link $(docker-compose ps -q postgres):postgres \
  auzom-migrations migrate -url \
  postgres://postgres:postgres@postgres/postgres?sslmode=disable up
