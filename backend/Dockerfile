FROM golang:1
RUN wget https://raw.githubusercontent.com/pote/gpm/master/bin/gpm
RUN chmod +x gpm
RUN mv gpm /usr/local/bin
WORKDIR /go/src/app
CMD gpm install 1>&2 && go install -v 1>&2 && tar -cf - $RUN -C /go/bin app
