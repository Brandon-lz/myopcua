# Build the application from source
FROM golang:1.21 AS build-stage
ENV GOPROXY=https://goproxy.cn
WORKDIR /app
COPY . .
RUN go build -o main .


# Deploy the application binary into a lean image
FROM debian AS build-release-stage
WORKDIR /app
COPY --from=build-stage ["/app/main", "./"]
COPY . .
RUN cp /app/health/healthcheck/healthcheck /healthcheck
EXPOSE 5001
# HEALTHCHECK  --interval=5m --timeout=3s CMD curl --fail http://localhost:6060/health || exit 1
# USER nonroot:nonroot
CMD ["./main"]
