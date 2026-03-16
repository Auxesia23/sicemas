FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o app cmd/api/*.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

ENV TZ=Asia/Jakarta

WORKDIR /app

COPY --from=builder /app/app .

RUN mkdir -p /app/data


ENTRYPOINT ["./app"]
