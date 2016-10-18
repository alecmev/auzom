#!/bin/bash

read -p "Rebuild migrations image? [yN] " build
if [ "$build" == "y" ]; then
  make migrations
fi

read -sp "Compose admin password: " PWD
printf "\n"
docker run --rm auzom-migrations migrate -url \
  $(echo $POSTGRES_URL_TEMPLATE | sed "s/{pwd}/$PWD/") up
