FROM golang:1.19.1-alpine3.16 AS builder
RUN apk add --no-cache git
WORKDIR /src
ADD . .
RUN go build -o /build/zetawallet ./cmd/zetawallet

FROM alpine:3.16
ENTRYPOINT ["zetawallet"]
RUN apk add --no-cache bash ca-certificates jq
COPY --from=builder /build/zetawallet /usr/local/bin/
