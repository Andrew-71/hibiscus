# syntax=docker/dockerfile:1

FROM golang:1.22 AS build-stage
WORKDIR /app

# Setup and compile Go
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /hibiscus

FROM alpine:3.19.1 AS deploy-stage
WORKDIR /

# Bring over the executable
COPY --from=build-stage /hibiscus /hibiscus

# Copy files
COPY public public/
COPY pages pages/
VOLUME data
VOLUME config

EXPOSE 7101
ENTRYPOINT ["/hibiscus"]