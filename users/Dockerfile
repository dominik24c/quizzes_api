FROM golang:1.19-alpine

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY ./internal ./internal
COPY ./users/core ./users/core
COPY ./users/main.go .

RUN go mod tidy
RUN go build main.go

EXPOSE 9999

CMD ["/app/main"]