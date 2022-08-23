FROM golang:1.19-alpine

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY ./internal ./internal
COPY ./solution/core ./solution/core
COPY ./solution/main.go .

RUN go mod tidy
RUN go build main.go

EXPOSE 9996

CMD ["/app/main"]