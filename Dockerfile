FROM golang:alpine
COPY . $GOPATH/src/url-shortener
WORKDIR $GOPATH/src/url-shortener
RUN apk add --no-cache git \
  && go get -d -v
RUN go build -o /go/bin/url-shortener
EXPOSE 8080
CMD ["/go/bin/url-shortener"]