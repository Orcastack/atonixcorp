FROM golang:1.27.1 AS build

# Put project inside GOPATH
WORKDIR /go/src/atonixcorp/cloud/database

COPY . .

RUN go mod tidy

RUN go build -o /bin/atdb-api ./cmd/atdb-api
RUN go build -o /bin/atdb-controller ./cmd/atdb-controller

FROM debian:12
LABEL org.opencontainers.image.source="https://github.com/orcadevstack/atonixcorp"

COPY --from=build /bin/atdb-api /usr/local/bin/atdb-api
COPY --from=build /bin/atdb-controller /usr/local/bin/atdb-controller

ENTRYPOINT ["/usr/local/bin/atdb-api"]
