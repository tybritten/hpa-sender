FROM golang:1.15


WORKDIR /go/src/app

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/hpa-sender .
#ENV GO111MODULE=on

FROM scratch

COPY --from=0 /app/hpa-sender /hpa-sender
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 8080


CMD ["/hpa-sender"]
