FROM golang:1.18

ADD . /satellite_code
WORKDIR /satellite_code

RUN set -ex; \
    apt-get update; \
    apt-get install -y netcat; \
    make linux; \
    mkdir -p /skywalking-satellite/configs; \
    cp configs/satellite_config.yaml /skywalking-satellite/configs/satellite_config.yaml; \
    cp bin/skywalking-satellite-latest-linux-amd64 /skywalking-satellite/satellite;

WORKDIR /skywalking-satellite

ENTRYPOINT ["./satellite", "start"]