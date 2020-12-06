FROM alpine:3.12

RUN apk add --no-cache ca-certificates

ADD ./grid-intensity-exporter /grid-intensity-exporter

ENTRYPOINT ["/grid-intensity-exporter"]
