FROM alpine:3.9
RUN apk update && apk add --no-cache ca-certificates
RUN mkdir /app
ADD . /app
WORKDIR /app
ENTRYPOINT [ "/app/cmd/switch-board-service/switch-board-service" ]
EXPOSE 7020
