#!/bin/bash

read -p "Rebuild migrations image? [yN] " build
if [ "$build" == "y" ]; then
  make migrations
fi

docker run --rm --network=host auzom-migrations migrate -path . -database \
  postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable up
