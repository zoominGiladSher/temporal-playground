FROM golang:1.22.4-alpine as worker_builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/worker/main.go

CMD ["/app/worker/main"]

