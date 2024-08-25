FROM golang:1.23-alpine AS build

ENV GOOS=linux

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd cmd/
COPY internal internal/

# RUN go test ./...
RUN go build -o /app/goyurback ./cmd/goyurback/main.go

FROM alpine:3.20

WORKDIR /app

RUN adduser -D goyurback && chown -R goyurback:goyurback /app
USER goyurback

COPY --chown=goyurback:goyurback --chmod=440 .env.* /app/
COPY --from=build --chown=goyurback:goyurback --chmod=770 /app/goyurback /app/

EXPOSE 8910

ENTRYPOINT ["./goyurback"]
