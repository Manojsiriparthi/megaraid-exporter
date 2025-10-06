FROM golang:1.19-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o megaraid-exporter .

FROM alpine:latest

# Install storcli
RUN apk --no-cache add curl unzip
RUN curl -L -o /tmp/storcli.zip https://docs.broadcom.com/docs-and-downloads/raid-controllers/raid-controllers-common-files/007.1914.0000.0000_Unified_StorCLI.zip
RUN cd /tmp && unzip storcli.zip && \
    unzip "Unified_StorCLI_*/Linux/storcli-*.noarch.rpm" && \
    tar xf storcli-*.tar.gz && \
    cp storcli /usr/sbin/storcli64 && \