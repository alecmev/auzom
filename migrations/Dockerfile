FROM golang:1
RUN go get -v github.com/mattes/migrate
WORKDIR /migrations
CMD tar -cf - $RUN *.sql -C /go/bin migrate
