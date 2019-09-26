FROM golang:1.13 AS build
WORKDIR /echoserver

ENV GOPROXY=https://proxy.golang.org
COPY go.mod /echoserver/
RUN go mod download

COPY main.go main.go
RUN CGO_ENABLED=0 GOOS=linux GOFLAGS=-ldflags=-w go build -o /go/bin/echo-server -ldflags=-s -v github.com/stevesloka/echo-server

FROM scratch AS final
COPY --from=build /go/bin/echo-server /bin/echo-server
