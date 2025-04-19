ARG GO_VERSION=1.24.2
FROM golang:${GO_VERSION} AS base
WORKDIR /src

COPY go.mod ./
RUN go mod download

FROM golang:${GO_VERSION} AS build
WORKDIR /src

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/server

FROM debian:bookworm-slim

ARG UID=1001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    appuser
USER appuser

COPY --from=build /bin/server /bin/

EXPOSE 8000

ENTRYPOINT ["/bin/server"]
