FROM alpine:3.15 AS builder

#install required ca-certs for gRPC
RUN  apk --no-cache --update add wget ca-certificates binutils
RUN mkdir -p /log


FROM scratch

WORKDIR app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY /artifacts/datadog-automation .
CMD ["./datadog-automation"]

