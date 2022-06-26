FROM golang:1.19

WORKDIR /go/src/app

COPY --from=cosmtrek/air /go/bin/air /usr/bin/air

CMD ["air", "--build.cmd", "go build -o bin/server cmd/server/main.go", "--build.bin", "bin/server"]