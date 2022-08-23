FROM golang:1.19-alpine

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY ./internal ./internal
COPY ./auth/core ./auth/core
COPY ./auth/main.go .

RUN go mod tidy
RUN go build main.go

EXPOSE 9998

CMD ["/app/main"]