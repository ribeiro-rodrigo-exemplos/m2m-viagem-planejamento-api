FROM alpine:3.7
RUN apk add --update bash
RUN apk add --no-cache tzdata
COPY m2m-viagem-planejamento-api.bin /
COPY ./configs/config.json /
WORKDIR /
ENTRYPOINT [ "./m2m-viagem-planejamento-api.bin" ]
CMD ["-config-location=./config.json"]