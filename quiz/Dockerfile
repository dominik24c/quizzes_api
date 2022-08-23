FROM golang:1.19-alpine

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY ./internal ./internal
COPY ./quiz/core ./quiz/core
COPY ./quiz/main.go .

RUN go mod tidy
RUN go build main.go

EXPOSE 9997

CMD ["/app/main"]