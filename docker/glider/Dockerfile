FROM golang:1.10.3-alpine3.8

# We want to ensure that release builds have cgo dependencies so we
# switch this on at the highest level.
ENV CGO_ENABLED 1
ENV GLIDE_HOME /.cache/go-build

RUN apk add --no-cache curl git bash build-base gcc musl-dev gmp  \
    # Install Glide (Go dependency manager)
    && curl https://glide.sh/get | sh \
    # Install Gox (Go cross-compiler)
    && go get -u -v github.com/mitchellh/gox \
    # Install ghr (Creates github releases and uploads artifacts)
    && go get -u -v github.com/tcnksm/ghr \
    && mkdir -p "$GLIDE_HOME" && chmod -R a+w "$GLIDE_HOME"
