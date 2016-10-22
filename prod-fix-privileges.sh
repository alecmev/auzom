#!/bin/bash

. .env.local
read -sp "Compose admin password: " PWD
printf "\n"
echo Fixing privileges...
docker run --rm postgres:$POSTGRES_VERSION psql \
  $(echo $POSTGRES_URL_TEMPLATE | sed "s/{pwd}/$PWD/") <<EOF
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO auzom;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO auzom;
EOF
echo Done
