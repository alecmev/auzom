#!/bin/bash

. .env
read -sp "Compose admin password: " PWD
printf "\n"
echo Backing up...
mkdir -p .backup/production
docker run --rm postgres:$POSTGRES_VERSION pg_dump --no-privileges --no-owner \
  --clean $(echo $POSTGRES_URL_TEMPLATE | sed "s/{pwd}/$PWD/") \
  > .backup/production/$(date -u +%Y-%m-%dT%H:%M:%SZ).sql
echo Done
