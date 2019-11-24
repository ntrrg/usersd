FROM golang:1.13.4-alpine3.10 AS build
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

