FROM golang:1.18 as build

ARG VERSION="latest"

ENV CGO_ENABLED=0

WORKDIR /src

COPY . .

RUN VERSION=$VERSION make linux
RUN mv /src/bin/skywalking-satellite-${VERSION}-linux-amd64 /src/bin/skywalking-satellite

FROM scratch

VOLUME /skywalking/configs

WORKDIR /skywalking

COPY --from=build /src/bin/skywalking-satellite /
COPY --from=build /src/configs /skywalking/configs

EXPOSE 12800 11800 1234

CMD ["/skywalking-satellite", "start", "--config", "/skywalking/configs/satellite_config.yaml"]
