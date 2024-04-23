FROM golang:latest AS build
COPY . /build
WORKDIR /build
ENV CGO_ENABLED=0
RUN go mod tidy
RUN go build -o thingnamer

FROM alpine:latest
WORKDIR /
COPY --from=build /build/thingnamer /thingnamer

ENTRYPOINT ["/thingnamer"]
