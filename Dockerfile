FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN --mount=type=cache,target=~/go/pkg/mod go mod download
COPY ./ ./
RUN --mount=type=cache,target=~/.cache/go-build go build -o ./altman ./cmd/bot/main.go
RUN ["chmod", "+x", "./altman"]
CMD ["/app/altman"]

FROM migrate/migrate AS migration
WORKDIR /migrations
COPY --from=builder /app/migrations/ /migrations

FROM alpine:latest AS prod
WORKDIR /app
COPY --from=builder /app/altman ./
RUN ["chmod", "+x", "./altman"]
CMD ["/app/altman"]