FROM golang:1.25 AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/reader-api ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/reader-migrate ./cmd/migrate
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/reader-dictionary-import ./cmd/dictionary-import
RUN mkdir -p /out/book-storage

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app
COPY --from=build /out/reader-api /app/reader-api
COPY --from=build /out/reader-migrate /app/reader-migrate
COPY --from=build /out/reader-dictionary-import /app/reader-dictionary-import
COPY --from=build --chown=nonroot:nonroot /out/book-storage/ /var/lib/english-reader/books/

EXPOSE 9090 8081
USER nonroot:nonroot
ENTRYPOINT ["/app/reader-api"]
