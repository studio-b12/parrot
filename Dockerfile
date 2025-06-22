FROM golang:alpine AS build
WORKDIR /build
COPY go.mod go.sum ./
COPY cmd/ cmd/
COPY pkg/ pkg/
RUN go build -o parrot cmd/parrot/main.go

FROM alpine:latest
COPY --from=build /build/parrot /usr/local/bin/parrot
EXPOSE 8080
ENTRYPOINT [ "/usr/local/bin/parrot" ]