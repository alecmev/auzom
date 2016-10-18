#!/bin/bash

if [ -z "$1" ]; then
  echo Tag name is mandatory, exiting
  exit 1
fi

run=Dockerfile-$2
if [ -z "$2" ]; then
  run=Dockerfile-run
fi
echo Using $run for runtime image

docker build --pull -t $1-builder .
docker run -v $(pwd):$(docker inspect -f '{{.Config.WorkingDir}}' $1-builder) \
  --rm -e RUN=$run $1-builder | docker build -f $run --pull --no-cache -t $1 -
