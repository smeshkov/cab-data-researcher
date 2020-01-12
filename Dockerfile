FROM alpine:3.4

RUN apk -U add ca-certificates

EXPOSE 8080

ADD ./_dist/cabresearcher_linux /bin/cabresearcher_linux
ADD _resources/config.yml /etc/cabresearcher/config/config.yml

CMD ["cabresearcher_linux", "-config", "/etc/cabresearcher/config/config.yml"]
