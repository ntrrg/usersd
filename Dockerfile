FROM golang:1.11-alpine3.8 AS build
RUN apk update && apk add git
RUN \
  cd /tmp && \
  wget -c 'https://github.com/magefile/mage/releases/download/v1.8.0/mage_1.8.0_Linux-64bit.tar.gz' && \
  tar -xf mage_1.8.0_Linux-64bit.tar.gz && \
  cp -af /tmp/mage $(go env GOPATH)/bin/
WORKDIR /src
COPY . .
RUN mage

FROM scratch
COPY --from=build /src/dist /bin
VOLUME "/data"
EXPOSE 4000
USER 1000
ENTRYPOINT ["/bin/usersd"]
CMD ["--db", "/data"]

