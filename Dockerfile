FROM golang:1.11-alpine3.8 AS build
RUN apk update && apk add git
WORKDIR /src
COPY . .
RUN ./mage

FROM scratch
COPY --from=build /src/dist /bin
VOLUME "/data"
EXPOSE 4000
USER 1000
ENTRYPOINT ["/bin/usersd"]
CMD ["--db", "/data"]

