FROM golang:1.22-alpine AS builder

WORKDIR /opt/app
ADD . .
RUN go build -o spaserver .

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
COPY --from=builder /opt/app/spaserver /usr/bin/spaserver

EXPOSE 80

ENTRYPOINT [ "/usr/bin/spaserver" ]

HEALTHCHECK --interval=10s --timeout=1s --start-period=5s --retries=3 \
	CMD [ "wget", "-qO", "-", "http://localhost/_health" ]