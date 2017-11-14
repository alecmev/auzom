FROM node:8

COPY . /app
RUN set -ex; \
  cd /app; \
  yarn --frozen-lockfile; \
  yarn run build

FROM nginx

COPY --from=0 /app/dist/* /data/
COPY --from=0 /app/nginx.conf /etc/nginx/
RUN chown -R nginx:nginx /data /etc/nginx/nginx.conf
