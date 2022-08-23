FROM golang:1.19-alpine

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY ./internal ./internal
COPY ./answers/core ./answers/core
COPY ./answers/main.go .

RUN go mod tidy
RUN go build main.go

EXPOSE 9995

CMD ["/app/main"]