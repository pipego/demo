FROM golang:latest AS build-stage
WORKDIR /go/src/app
COPY . .
RUN apt update && \
    apt install -y upx
RUN make build

FROM gcr.io/distroless/base-debian11 AS production-stage
WORKDIR /
COPY --from=build-stage /go/src/app/bin/demo /
COPY --from=build-stage /go/src/app/config/config.yml /
USER nonroot:nonroot
CMD ["/demo", "--config-file=/config.yml"]
