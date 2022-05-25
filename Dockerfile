FROM golang:latest AS build-stage
WORKDIR /go/src/app
COPY . .
RUN apt update && \
    apt install -y upx
RUN make build

FROM gcr.io/distroless/base-debian11 AS production-stage
WORKDIR /
COPY --from=build-stage /go/src/app/bin/demo /
COPY --from=build-stage /go/src/app/test/config/config.yml /
COPY --from=build-stage /go/src/app/test/data/runner.json /
COPY --from=build-stage /go/src/app/test/data/scheduler.json /
USER nonroot:nonroot
CMD ["/demo", "--config-file=/config.yml", "--runner-file=/runner.json", "--scheduler-file=/scheduler.json"]
