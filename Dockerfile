FROM golang:1.26.8-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal

ENV CGO_ENABLED=0 GOOS=linux
RUN go build -ldflags="-s -w" -o /out/server    ./cmd/server && \
  go build -ldflags="-s -w" -o /out/migration ./cmd/migration && \
  go build -ldflags="-s -w" -o /out/seeder    ./cmd/seeder

RUN mkdir -p /app/uploads

FROM gcr.io/distroless/static-debian12
COPY --from=build /out/ /
WORKDIR /app
COPY --from=build --chown=65532:65532 /app/uploads /app/uploads
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/server"]
