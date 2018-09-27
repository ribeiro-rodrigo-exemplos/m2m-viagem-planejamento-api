# build stage
FROM golang:1.11.0-alpine3.7 AS build-env
RUN apk add --update --no-cache git
ADD . /src
ADD ./go-logging-package-level /go-logging-package-level
RUN cd /src && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o m2m-viagem-planejamento-api.bin cmd/viagemPlanejamentoAPI/main.go

FROM alpine:3.7
RUN apk add --update bash
RUN apk add --no-cache tzdata
COPY --from=build-env /src/m2m-viagem-planejamento-api.bin /
COPY ./configs/config.json /
WORKDIR /
ENTRYPOINT [ "./m2m-viagem-planejamento-api.bin" ]
CMD ["-config-location=./config.json"]