ARG BASE_IMAGE=docker.io/ubuntu:24.04
ARG BUILD_IMAGE=docker.io/library/golang:1.24

FROM ${BUILD_IMAGE} as builder
WORKDIR /zenlayer-cloud-csi-driver
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -mod=vendor  -ldflags "-s -w" -o  _output/zeccsi ./cmd/

FROM ${BASE_IMAGE} 
LABEL maintainers="Zenlayer"
LABEL description="ZEC CSI plugin"

RUN apt-get update
RUN apt-get install -y util-linux e2fsprogs xfsprogs mount ca-certificates udev && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /zenlayer-cloud-csi-driver/_output/zeccsi /zeccsi
RUN chmod +x /zeccsi && \
    mkdir -p /var/log/zeccsi
ENTRYPOINT ["/zeccsi"]