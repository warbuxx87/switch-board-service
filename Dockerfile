FROM alpine:3.9
RUN apk update && apk add --no-cache ca-certificates
RUN mkdir /app
COPY . /app
WORKDIR /app
EXPOSE 80
ENTRYPOINT [ "/app/cmd/switch-board-service/switch-board-service", "prod"]
