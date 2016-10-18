#!/bin/bash

docker exec -i $(docker-compose ps -q postgres) \
  psql postgres postgres <<EOF
DROP SCHEMA public CASCADE;
CREATE SCHEMA public;
EOF

docker exec -i $(docker-compose ps -q postgres) psql postgres postgres < \
  $(command ls -d .backup/production/* | tail -1)
