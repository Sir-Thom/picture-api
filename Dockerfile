FROM  docker.io/arm64v8/golang:1.21-alpine
WORKDIR /app
COPY . .
RUN go build -o picture-api .
EXPOSE 8080
RUN apk update && \
    apk add --no-cache curl
CMD ["./picture-api"]

