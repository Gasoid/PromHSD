FROM golang:1.19-bullseye AS builder
WORKDIR /code
ADD go.mod /code/
#ADD go.sum /code/
RUN go mod download
ADD . /code/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /code/promhsd .
RUN chmod a+x /code/promhsd


FROM alpine:3.6
WORKDIR /root/
RUN apk --no-cache --update add bash curl less jq openssl
COPY --from=builder /code/promhsd /usr/local/bin/promhsd
CMD exec /usr/local/bin/promhsd
