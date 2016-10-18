FROM node:6
WORKDIR /app
CMD npm run build 1>&2 && tar -cf - $RUN nginx.conf dist/
