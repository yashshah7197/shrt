# Build the Go binary
FROM golang:1.17 as build_shrt
ENV CGO_ENABLED 0
ARG BUILD_REF

# Copy the source code in to the container
COPY . /shrt

# Build the shrt binary
WORKDIR /shrt
RUN go build -ldflags "-X main.build=${BUILD_REF}"

# Run the shrt binary in Alpine
FROM alpine:3.15
ARG BUILD_DATE
ARG BUILD_REF
COPY --from=build_shrt /shrt/shrt /shrt/shrt
WORKDIR /shrt
CMD ["./shrt"]

LABEL org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.title="shrt" \
      org.opencontainers.image.authors="Yash Shah <yash@yashshah.io>" \
      org.opencontainers.image.source="https://github.com/yashshah7197/shrt" \
      org.opencontainers.image.revision="${BUILD_REF}" \
      org.opencontainers.image.vendor="Yash Shah"
