FROM node:alpine AS node-builder

WORKDIR /build

COPY package.json package-lock.json ./

RUN npm install

COPY tsconfig.base.json vite.config.ts ./

COPY ./styles ./styles
COPY ./scripts ./scripts
COPY ./react ./react
COPY ./templates ./templates

RUN npx tailwindcss -i ./styles/style.css -o ./static/style.css -m
RUN npx tsc -p scripts
RUN npx vite build


FROM golang:alpine AS go-builder

RUN apk update && apk add --no-cache \
    gcc \
    musl-dev

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY ./cmd ./cmd
COPY ./internal ./internal
COPY main.go ./

ENV CGO_ENABLED=1

RUN go build -o blog .


FROM alpine

WORKDIR /app

COPY ./locales ./locales
COPY ./templates ./templates
COPY ./static ./static
COPY ./migrations ./migrations
COPY config.json ./

RUN mkdir ./database
RUN mkdir ./uploads

COPY --from=node-builder /build/static ./static
COPY --from=go-builder /build/blog .

EXPOSE 8000

CMD ["./blog", "-p", "8000", "-e", "prod"]
