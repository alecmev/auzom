#!/bin/bash

# stackoverflow.com/questions/7454526
EXEC=( "-c" "pkill gulp; pkill app; pkill node; pkill npm; bash" )
if [ "$2" = "nokill" ]; then EXEC=( ); fi

docker exec -it $(docker-compose ps -q $1) bash "${EXEC[@]}"
