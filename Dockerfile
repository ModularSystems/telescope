FROM golang:alpine AS builder

RUN apk update && apk add --no-cache \
 make \
 gcc \
 libc-dev
WORKDIR $GOPATH/src/github.com/modularsystems/telescope
COPY . .
# Fetch dependencies w/ go mod
RUN go mod download
# Build the binary.
RUN go build -o /go/bin/telescope ./cmd/telescope/main.go

####################################################################

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