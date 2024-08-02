FROM golang:1.22-alpine AS builder

WORKDIR /opt/app
ADD . .
RUN go build -o spaserver .

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
COPY --from=builder /opt/app/spaserver /usr/bin/spaserver

ENTRYPOINT [ "/usr/bin/spaserver" ]