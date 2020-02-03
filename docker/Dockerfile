FROM golang:1.13 as builder

WORKDIR /go/src/github.com/linden-honey/linden-honey-scraper-go

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY pkg ./pkg
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go install -v -ldflags="-w -s" ./cmd/server

FROM scratch

LABEL \
    name="linden-honey-scraper-go" \
    maintainer="aliaksandr.babai@gmail.com"

ARG WORK_DIR=/linden-honey
WORKDIR $WORK_DIR

ENV SERVER_PORT=8080
EXPOSE $SERVER_PORT

COPY --from=builder /go/bin/server /bin/server
COPY api/ ./api

ENTRYPOINT [ "/bin/server" ]
