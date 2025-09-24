FROM golang:1.25-alpine

WORKDIR /app

RUN apk add --no-cache make

COPY cmd cmd
COPY db db
COPY internal internal
COPY pkg pkg
COPY go.mod go.mod
COPY go.sum go.sum
COPY Makefile Makefile
COPY config.yaml config.yaml

RUN go mod download
RUN make build

EXPOSE 8080

CMD ["./bin/run"]