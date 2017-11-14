FROM golang:1

COPY . /go/src/app
RUN set -ex; \
  cd /go/src/app; \
  wget https://raw.githubusercontent.com/pote/gpm/master/bin/gpm; \
  chmod +x gpm; \
  ./gpm install; \
  go install

FROM frolvlad/alpine-glibc

RUN apk --update add ca-certificates
COPY --from=0 /go/bin/app /usr/local/bin/
EXPOSE 80
ENTRYPOINT ["app"]
