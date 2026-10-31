# syntax=docker/dockerfile:1
# Multi-stage, multi-arch build (works on the Pi and on x86).

# ---- build stage: cross-compiles a static binary for the target platform
FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS builder
ARG TARGETOS
ARG TARGETARCH
WORKDIR /src

# Download modules first so this layer is cached between code changes
COPY go.mod go.sum ./
RUN go mod download

COPY . .
# CGO off = fully static binary: no glibc/musl loader needed at runtime.
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags="-s -w" -o /out/picture-api .

# ---- runtime stage: only the binary (+ curl for the compose healthcheck)
FROM alpine:3.22
RUN apk add --no-cache curl ca-certificates tzdata
WORKDIR /app
COPY --from=builder /out/picture-api /app/picture-api
EXPOSE 8080
CMD ["/app/picture-api"]