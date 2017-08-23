FROM node:8
WORKDIR /app
CMD yarn run build 1>&2 && tar -cf - $RUN nginx.conf dist/
