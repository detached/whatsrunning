FROM golang:latest
RUN mkdir -p /opt/whatsrunning/

COPY whatsrunning /opt/whatsrunning/
COPY templates /opt/whatsrunning/
COPY dist /opt/whatsrunning/
COPY config.json /opt/whatsrunning/

EXPOSE 8080

WORKDIR /opt/whatsrunning/

ENTRYPOINT ["./whatsrunning"]
