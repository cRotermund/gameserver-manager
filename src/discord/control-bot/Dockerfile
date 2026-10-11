FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bot .

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /root/

COPY --from=builder /app/bot .

CMD ["./bot"]