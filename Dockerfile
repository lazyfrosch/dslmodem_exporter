FROM alpine
COPY dslmodem_exporter /usr/bin/
ENTRYPOINT ["/usr/bin/dslmodem_exporter"]
