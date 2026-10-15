FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY src/discord/control-bot ./src/discord/control-bot
COPY src/pkg/controlplaneclient ./src/pkg/controlplaneclient
RUN go work init ./src/discord/control-bot ./src/pkg/controlplaneclient
RUN go work edit -replace=github.com/cRotermund/gameserver-manager/src/pkg/controlplaneclient@v0.0.0=./src/pkg/controlplaneclient
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bot ./src/discord/control-bot

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /root/

COPY --from=builder /app/bot .

CMD ["./bot"]
