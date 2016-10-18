#!/bin/bash

echo Initializing...
LOCAL_IP=`hostname -I | awk '{print $1}'`

NGINX_ENV=nginx/.env
echo Generating $NGINX_ENV...
cat > $NGINX_ENV <<EOF
VIRTUAL_HOST=$LOCAL_IP
LETSENCRYPT_HOST=$LOCAL_IP
LETSENCRYPT_EMAIL=admin@$LOCAL_IP
EOF

BACKEND_ENV=backend/.env
echo Generating $BACKEND_ENV...
cat > $BACKEND_ENV <<EOF
STATIC_HOST=$LOCAL_IP
EOF

ROOT_ENV=.env
echo Generating $BACKEND_ENV...
cat > $ROOT_ENV <<EOF
POSTGRES_VERSION=9.5.4
POSTGRES_URL_TEMPLATE=postgres://admin:{pwd}@example.com:12345/auzom
DOCKER_REGISTRY=example.dkr.ecr.antarctica.amazonaws.com
EOF

echo Done
