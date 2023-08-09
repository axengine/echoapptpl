# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:1.19 as build

# Set go proxy
ENV GOPROXY="https://goproxy.cn"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN make build-linux

# Start from a new base image as runtime environment
FROM alpine:3.9

LABEL maintainer="axengine<axengine@gmail.com>"

# Set locale
ENV LANG en_US.UTF-8
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

ENV TIMEZONE "Asia/Shanghai"

RUN apk add --no-cache tzdata &&\
	cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime &&\
	echo $TIMEZONE >  /etc/timezone &&\
	apk del tzdata

WORKDIR /app
COPY --from=build /app/build/echoapp /app/main

ENTRYPOINT ["/app/main"]
