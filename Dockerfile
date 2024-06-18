FROM  docker.io/arm64v8/golang:1.21-alpine
WORKDIR /app
COPY . .
RUN go build -o picture-api .
EXPOSE 8080
CMD ["./picture-api"]

