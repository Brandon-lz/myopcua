# Build the application from source
FROM golang:1.22 AS build-stage
ENV GOPROXY=https://goproxy.cn
WORKDIR /app
COPY . .
RUN go build -o main .
WORKDIR /app/health/healthcheck
RUN go build -o healthcheck main.go


# Deploy the application binary into a lean image
FROM debian AS build-release-stage
WORKDIR /app
COPY --from=build-stage ["/app/main","/app/health/healthcheck/healthcheck", "./"]
COPY . .
RUN cp /app/healthcheck /healthcheck
EXPOSE 8080
# HEALTHCHECK  --interval=5m --timeout=3s CMD curl --fail http://localhost:6060/health || exit 1
# USER nonroot:nonroot
CMD ["./main"]
