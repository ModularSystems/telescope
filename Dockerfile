############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder
# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/github.com/modularsystems/telescope

COPY . .
# Fetch dependencies w/ go mod
RUN go mod download
# Build the binary.
RUN make test
RUN make build

############################
# STEP 2 build a small image w/ Ruby and our wpscan binary
############################
FROM alpine:3

RUN apk -U add \
    alpine-sdk \
    ruby \
    ruby-dev \
    libffi-dev \
    zlib-dev
RUN gem install wpscan
RUN apk del alpine-sdk

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