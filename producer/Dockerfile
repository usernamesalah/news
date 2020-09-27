# builder
FROM golang:1.14-alpine as go-builder

RUN apk update && apk upgrade && \
    apk --update add git make

WORKDIR /src

COPY . .

RUN go build -o bin/producer .

# distribution
FROM alpine:3.12

COPY --from=go-builder /src/bin/producer .

COPY migrations migrations

CMD ["./producer"]
