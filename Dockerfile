FROM golang:1.11-alpine3.8 AS build
WORKDIR /go/src/github.com/ntrrg/usersd
COPY vendor vendor
COPY pkg pkg
COPY api api
COPY main.go .
RUN CGO_ENABLED=0 go build -o $(go env GOPATH)/bin/usersd

FROM alpine3.8 as debug
COPY --from=build /go/bin /bin

FROM scratch
COPY --from=build /go/bin /bin
VOLUME "/data"
EXPOSE 4000
ENTRYPOINT ["/bin/usersd", "-db", "/data"]

