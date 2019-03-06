FROM golang:latest AS build

WORKDIR /geo-api
COPY . .

RUN CGO_ENABLED=0 go build -ldflags "-s -w" -a -installsuffix cgo -o build/geo-api


FROM alpine:latest AS prod

WORKDIR /root/
COPY --from=build /geo-api/build/geo-api .
COPY schema/geo-schema.json schema/geo-schema.json

ENTRYPOINT ["./geo-api"]
