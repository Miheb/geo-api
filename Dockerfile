FROM golang:latest AS build

WORKDIR /go/src/github.com/campus-iot/geo-api
COPY . .

RUN curl -fsSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 && \
    chmod +x /usr/local/bin/dep && \
    dep ensure -vendor-only && \
    CGO_ENABLED=0 go build -ldflags "-s -w" -a -installsuffix cgo -o build/geo-api


FROM alpine:latest AS prod

WORKDIR /root/

COPY --from=build /go/src/github.com/campus-iot/geo-api/build/geo-api .
COPY schema/geo-schema.json schema/geo-schema.json

ENTRYPOINT ["./geo-api"]
