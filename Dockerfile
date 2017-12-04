FROM golang:1.8-alpine as builder
RUN apk --no-cache add make git
RUN mkdir -p /go/src/github.com/garethr/kubetest/
COPY . /go/src/github.com/garethr/kubetest/
WORKDIR /go/src/github.com/garethr/kubetest/
RUN make linux

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/src/github.com/garethr/kubetest/bin/linux/amd64/kubetest .
ENTRYPOINT ["/kubetest"]
CMD ["--help"]
