FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
RUN go mod tidy

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -v -o /server ./cmd/server

FROM alpine:latest

WORKDIR /

COPY --from=builder /server /server
COPY --from=builder /app/queries ./queries

EXPOSE 8080

ENTRYPOINT [ "/server" ]
