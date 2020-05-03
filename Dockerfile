FROM alpine:latest as build
RUN apk --update add ca-certificates
FROM scratch
ENV PATH=/bin
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ADD public /public
ADD gowdc /
CMD ["./gowdc"]
