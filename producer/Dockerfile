FROM golang:1.14-alpine as go-builder
WORKDIR /src
COPY . .
RUN go build -o bin/soccer .

FROM alpine:3.12
COPY --from=go-builder /src/bin/soccer .
COPY migrations migrations

CMD ["./soccer"]
