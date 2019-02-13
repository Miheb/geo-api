FROM golang:1.11-alpine AS development

WORKDIR /go/src/github.com/campus-iot/geo-API
COPY . .

RUN apk add --no-cache ca-certificates git

RUN go get -u ./... && \
    go build


FROM alpine:latest AS production

WORKDIR /root/
COPY --from=development /go/src/github.com/campus-iot/geo-API/geo-API .
COPY test/data.json test/data.json
COPY schema/geo-schema.json schema/geo-schema.json

CMD ["./geo-API"]
