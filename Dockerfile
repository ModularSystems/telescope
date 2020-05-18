FROM golang:alpine AS builder

# BUG - https://github.com/golang/go/issues/27303
ENV CGO_ENABLED 0

WORKDIR $GOPATH/src/github.com/modularsystems/telescope
COPY . .
# Fetch dependencies w/ go mod
RUN go mod download
# catch test failures on build
RUN go test -v ./...
# Build the binary.
RUN go build -o /go/bin/telescope ./cmd/telescope/main.go

####################################################################

FROM ruby:2.6.3-alpine

# You can see the requirements for wpscan in their upstream Dockerfile
# https://github.com/wpscanteam/wpscan/blob/master/Dockerfile
RUN apk -U add \
    git \
    gcc \
    libcurl \
    libffi-dev \
    libxml2 \
    make \
    musl-dev \
    ruby \
    ruby-dev \
    procps \
    sqlite-dev \
    sqlite-libs \
    zlib-dev
RUN gem install --verbose --backtrace --debug wpscan

# Copy our static executable.
COPY --from=builder /go/bin/telescope /go/bin/telescope

# Configuration is done through the following environment variables
# If empty, we default to the minimal default behavior possible
ENV SENDGRID_API_KEY=""
ENV SENDGRID_SENDER_NAME=""
ENV SENDGRID_SENDER_EMAIL=""
ENV WPVULNDB_API_KEY=""

# Run the telescope binary.
ENTRYPOINT ["/go/bin/telescope"]