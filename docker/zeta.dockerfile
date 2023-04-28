FROM golang:1.19.1-alpine3.16 AS builder
RUN apk add --no-cache git
ENV CGO_ENABLED=0
WORKDIR /src
ADD . .
RUN go build -o /build/zeta ./cmd/zeta

FROM alpine:3.16
ENTRYPOINT ["zeta"]
RUN apk add --no-cache bash ca-certificates jq
COPY --from=builder /build/zeta /usr/local/bin/
