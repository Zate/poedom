FROM golang:latest as builder
WORKDIR /go/src/github.com/zate/poedom/
RUN go get -u github.com/golang/dep/cmd/dep
COPY main.go .
COPY favicon.ico .
COPY common.css .
COPY public public
COPY static static
RUN dep init && dep ensure
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o poedom .

FROM scratch
WORKDIR /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /go/src/github.com/zate/poedom/poedom .
COPY --from=builder /go/src/github.com/zate/poedom/common.css .
COPY --from=builder /go/src/github.com/zate/poedom/favicon.ico .
COPY --from=builder /go/src/github.com/zate/poedom/public /public
COPY --from=builder /go/src/github.com/zate/poedom/static /static
EXPOSE 2086
CMD ["/poedom"]