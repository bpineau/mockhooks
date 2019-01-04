FROM golang:1.11.4-stretch as builder
WORKDIR /root
COPY . .
RUN CGO_ENABLED=0 go build

FROM alpine:3.8
COPY --from=builder /root/mockhooks /usr/bin/
RUN apk add tini
ENTRYPOINT ["/sbin/tini", "--"]
USER nobody
CMD ["/usr/bin/mockhooks"]
