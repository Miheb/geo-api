FROM golang:latest

COPY ./swagger $GOPATH/src/swagger

WORKDIR $GOPATH/src/timo69/main

COPY . .

RUN go get ./...
RUN go install -v ./...

EXPOSE 8081

CMD ["main"]
