FROM golang:1.18.2-alpine3.16 AS builder

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . ./
RUN go build -o /binary

FROM alpine:3.16

WORKDIR /app
COPY --from=builder /binary ./

CMD [ "./binary" ]