FROM golang:1.11 AS build
WORKDIR /go/src/github.com/campus-iot/geo-API
COPY . .
RUN go get -d -v ./... && go build

FROM alpine:latest AS prod
WORKDIR /root/
COPY --from=build /go/src/github.com/campus-iot/geo-API/geo-API .
COPY test/data.json test/data.json
COPY schema/geo-schema.json schema/geo-schema.json

CMD ["./geo-API"]
