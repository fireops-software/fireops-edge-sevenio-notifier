FROM alpine:latest AS build-stage

# install golang
WORKDIR /
RUN GO_VERSION=1.24.2 \
    && wget https://go.dev/dl/go$GO_VERSION.linux-amd64.tar.gz \
    && tar -xzf go$GO_VERSION.linux-amd64.tar.gz \
    && rm go$GO_VERSION.linux-amd64.tar.gz
ENV PATH=$PATH:/go/bin

# set workdir for project
WORKDIR /app
COPY . .
RUN go build -o fireops-edge-sevenio-notifier main.go

# Deploy the application binary into a lean image
FROM alpine:latest AS build-release-stage

WORKDIR /app

COPY --from=build-stage /app/fireops-edge-sevenio-notifier /app

ENTRYPOINT ["/app/fireops-edge-sevenio-notifier"]
