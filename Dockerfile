FROM golang:1.15


WORKDIR /go/src/app

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/ce-webhook-sender .
#ENV GO111MODULE=on

FROM scratch

COPY --from=0 /app/ce-webhook-sender /ce-webhook-sender

EXPOSE 8080


CMD ["/ce-webhook-sender"]
