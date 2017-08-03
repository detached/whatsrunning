FROM golang:latest
RUN mkdir -p /opt/whatsrunning/data

COPY whatsrunning /opt/whatsrunning/
COPY templates /opt/whatsrunning/
COPY dist /opt/whatsrunning/
COPY config.json /opt/whatsrunning/

EXPOSE 8080

WORKDIR /opt/whatsrunning/
VOLUME /opt/whatsrunning/data

ENTRYPOINT ["./whatsrunning"]
