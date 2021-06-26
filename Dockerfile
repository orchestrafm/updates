# Set Builder Image
FROM golang:1.12.17-alpine3.10 as builder

# Add Build Dependencies and Working Directory
RUN apk --no-cache add build-base git tar wget
RUN wget https://github.com/go-task/task/releases/download/v2.0.0/task_linux_amd64.tar.gz
RUN tar zxvf task_linux_amd64.tar.gz && cp ./task /usr/bin/
RUN chmod +x /usr/bin/task
RUN mkdir /build
ADD . /build/
WORKDIR /build

# Compile
ENV GO111MODULE=on
RUN go install github.com/gobuffalo/packr/packr
RUN task build OUTNAME=service

# Move to Base Image and Run
FROM alpine:3.12.0
RUN apk update && apk upgrade && apk add ca-certificates
COPY --from=builder /build/service /app/
WORKDIR /app
CMD ["./service"]