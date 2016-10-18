#!/bin/bash

docker kill nginx-letsencrypt
docker rm nginx-letsencrypt
docker kill nginx-proxy
docker rm nginx-proxy

docker pull jwilder/nginx-proxy
docker run -d --name nginx-proxy -p 80:80 -p 443:443 \
  -v /etc/nginx/certs:/etc/nginx/certs:ro \
  -v /etc/nginx/vhost.d \
  -v /usr/share/nginx/html \
  -v /var/run/docker.sock:/tmp/docker.sock:ro \
  jwilder/nginx-proxy

docker pull jrcs/letsencrypt-nginx-proxy-companion
docker run -d --name nginx-letsencrypt \
  -v /etc/nginx/certs:/etc/nginx/certs \
  --volumes-from nginx-proxy \
  -v /var/run/docker.sock:/var/run/docker.sock:ro \
  jrcs/letsencrypt-nginx-proxy-companion
