#!/bin/bash

echo Initializing...

BACKEND_ENV=backend/.env.local
echo Generating $BACKEND_ENV...
cat > $BACKEND_ENV <<EOF
STATIC_HOST=localhost:3000
EOF

ROOT_ENV=.env.local
echo Generating $ROOT_ENV...
cat > $ROOT_ENV <<EOF
. .env
POSTGRES_URL_TEMPLATE=postgres://admin:{pwd}@example.com:12345/auzom
DOCKER_REGISTRY=example.dkr.ecr.antarctica.amazonaws.com
EOF

echo Done
