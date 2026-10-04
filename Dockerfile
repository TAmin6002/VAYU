# ---- build stage ----
FROM golang:1.22-alpine AS build
WORKDIR /src

# If proxy.golang.org is unreachable from your network, build with e.g.
#   docker compose build --build-arg GOPROXY=https://goproxy.io,direct
ARG GOPROXY=https://proxy.golang.org,direct
ENV GOPROXY=${GOPROXY}

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -trimpath -o /out/server ./cmd/server

# ---- runtime stage ----
FROM alpine:3.20
WORKDIR /app
RUN adduser -D -u 10001 appuser

COPY --from=build /out/server /app/server
COPY web /app/web

USER appuser
EXPOSE 8080
CMD ["/app/server"]
