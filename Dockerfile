FROM golang:1.16-alpine3.13 AS build
RUN apk update && apk add make
WORKDIR /src
COPY . .
RUN make

FROM scratch
COPY --from=build /src/dist/usersd-linux-amd64 /bin/usersd
USER 1000
EXPOSE 4000
VOLUME "/var/usersd"
ENTRYPOINT ["/bin/usersd", "--db", "/var/usersd"]

