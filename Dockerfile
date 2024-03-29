FROM golang:1.21-alpine AS builder
ENV GO111MODULE=on
RUN apk --update upgrade \
    && apk --no-cache --no-progress add git
WORKDIR /src
ADD . .
RUN go mod download
RUN go mod verify
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o dnsu ./cmd/dnsu

FROM alpine:3.18
LABEL maintainer="Levent SAGIROGLU <LSagiroglu@gmail.com>"
VOLUME /dnsu
COPY --from=builder /src/dnsu /dnsu/dnsu
COPY --from=builder /src/update-dns.sh /dnsu/update-dns.sh
ENTRYPOINT ["tail"]
CMD ["-f","/dev/null"]

# docker build -t netyazilim/dnsu:0.7.2 .