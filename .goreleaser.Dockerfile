FROM alpine:3 AS builder

RUN apk update && apk upgrade && apk add --no-cache ca-certificates
WORKDIR /app
RUN echo "nobody:x:65534:65534:Nobody:/:" > /app/passwd

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/passwd /etc/passwd
COPY ./oci-resolve /bin/oci-resolve
USER nobody
WORKDIR /app
CMD ["/bin/oci-resolve"]
