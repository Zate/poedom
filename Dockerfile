FROM golang:latest as builder
WORKDIR /go/src/github.com/zate/poedom/
RUN go get -u github.com/golang/dep/cmd/dep
ADD main.go .
ADD public .
ADD static .
COPY public/* public/
COPY static/* static/
RUN dep init && dep ensure
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o poedom .

FROM scratch
WORKDIR /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /go/src/github.com/zate/poedom/poedom .
EXPOSE 2086
CMD ["/poedom"]