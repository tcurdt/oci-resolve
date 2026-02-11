FROM golang:1.26-alpine AS builder

RUN apk update && apk upgrade && apk add --no-cache ca-certificates
ADD . /app
RUN echo "nobody:x:65534:65534:Nobody:/:" > /app/passwd
WORKDIR /app
RUN ls -la
RUN go build ./...

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/passwd /etc/passwd
COPY --from=builder /app/oci-resolve /usr/local/bin/oci-resolve
USER nobody
WORKDIR /
CMD ["/usr/local/bin/oci-resolve"]
