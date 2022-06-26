FROM golang:1.18.3

WORKDIR /go/src/app

COPY --from=cosmtrek/air /go/bin/air /usr/bin/air

CMD ["air", "--build.cmd", "go build -o tmp/main cmd/server/main.go"]