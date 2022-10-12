FROM golang:1.16-alpine3.13 as builder

WORKDIR /go/src/github.com/linden-honey/linden-honey-scraper-go

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY pkg ./pkg
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go install -v -ldflags="-w -s" ./cmd/...

FROM alpine:3.13

COPY --from=builder /go/bin/server /bin/server

ENV SERVER_HOST=0.0.0.0 \
    SERVER_PORT=80
EXPOSE $SERVER_PORT

ENTRYPOINT [ "/bin/server" ]
