FROM golang:1.21-alpine3.19 as builder

WORKDIR /go/src/github.com/linden-honey/linden-honey-scraper-go

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY pkg ./pkg
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go install -v -ldflags="-w -s" ./cmd/...

FROM alpine:3.19

COPY --from=builder /go/bin/server /bin/scraper

ENTRYPOINT [ "/bin/scraper" ]
