FROM alpine:3.4

RUN apk -U add ca-certificates

EXPOSE 8080

ADD cab-data-researcher /bin/cab-data-researcher
ADD _resources/config.yml /etc/cab-data-researcher/config/config.yml

CMD ["cab-data-researcher", "-config", "/etc/cab-data-researcher/config/config.yml"]
