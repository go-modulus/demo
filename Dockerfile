FROM golang:1.19-alpine as tools

RUN apk add --no-cache build-base

RUN go install github.com/cespare/reflex@v0.3
RUN go install -tags 'nowasm' github.com/kyleconroy/sqlc/cmd/sqlc@v1.16
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.15

FROM golang:1.19-alpine

WORKDIR /go/src/app

RUN apk add --no-cache make build-base
COPY --from=tools /go/bin/reflex /usr/bin/reflex
COPY --from=tools /go/bin/sqlc /usr/bin/sqlc
COPY --from=tools /go/bin/migrate /usr/bin/migrate

CMD ["make", "start"]