FROM ubuntu:20.10

LABEL org.opencontainers.image.source https://github.com/ovrclk/chain-supply

EXPOSE 8080

RUN apt-get update && \
    apt-get install -y ca-certificates && \
    update-ca-certificates

ADD chain-supply chain-supply

ENTRYPOINT ["/chain-supply"]

CMD ["server"]
