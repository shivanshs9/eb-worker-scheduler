FROM golang:1.18-alpine

RUN apk add --update --no-cache dumb-init git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o ./crond

ENTRYPOINT [ "/usr/bin/dumb-init", "--", "/app/crond" ]
