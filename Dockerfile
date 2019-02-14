FROM golang:1.11 AS build
WORKDIR /go/src/github.com/campus-iot/geo-api
COPY . .
RUN go get -d -v ./... && go build

FROM alpine:latest AS prod
WORKDIR /root/
COPY --from=build /go/src/github.com/campus-iot/geo-api/geo-api .
COPY schema/geo-schema.json schema/geo-schema.json

CMD ["./geo-api"]
