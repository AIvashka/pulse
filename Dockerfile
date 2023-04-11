FROM golang:1.19-alpine3.15 AS build

RUN apk add --no-cache git

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main app/app.go
FROM alpine:3.15
WORKDIR /app
COPY --from=build /app/main ./
COPY config.json ./
COPY config.yaml ./

CMD ["./main"]


#docker run -d -v "$(pwd)/.env:/app/.env" --restart=always --name spread_recorder app


#docker build -t app  .
#docker run -d -v "$(pwd)/.env:/app/.env" --name spread_recorder app