FROM golang:1.22.4-alpine as gateway_builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN ls -la

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/gateway/main.go

EXPOSE 8091

CMD ["/app/gateway/main"]

