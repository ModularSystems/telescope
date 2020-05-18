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

FROM alpine:3

RUN apk -U add \
    alpine-sdk \
    ruby \
    ruby-dev \
    libffi-dev \
    zlib-dev
RUN gem install wpscan
RUN apk del alpine-sdk ruby-dev libffi-dev zlib-dev

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