FROM golang:1.22-alpine AS builder
RUN apk update && apk upgrade \
 && apk add --no-cache ca-certificates
WORKDIR /app
ADD . /app
RUN echo "nobody:x:65534:65534:Nobody:/:" > /app/passwd
RUN go build ./...

FROM scratch
WORKDIR /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/passwd /etc/passwd
COPY --from=builder /app/oci-resolve /bin/oci-resolve
USER nobody
CMD ["/bin/oci-resolve"]

# HEALTHCHECK CMD curl --fail http://localhost/api/healthcheck || exit 1
# HEALTHCHECK NONE
# HEALTHCHECK CMD exit 0