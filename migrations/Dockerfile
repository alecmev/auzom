FROM golang:1
RUN go get -u -d github.com/mattes/migrate/cli github.com/lib/pq
RUN go build -tags 'postgres' -o /usr/local/bin/migrate github.com/mattes/migrate/cli
WORKDIR /migrations
CMD tar -cf - $RUN *.sql -C /usr/local/bin migrate
