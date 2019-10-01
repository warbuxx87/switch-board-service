FROM alpine:3.9
RUN apk update && apk add --no-cache ca-certificates
RUN mkdir /app
ADD . /app
WORKDIR /app/cmd/switch-board-service
RUN ls
ENTRYPOINT [ "./switch-board-service" ]