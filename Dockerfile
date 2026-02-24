FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod .
RUN go mod download || true
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/zenx ./cmd/zenx

FROM gcr.io/distroless/static-debian12
COPY --from=build /out/zenx /usr/local/bin/zenx
ENTRYPOINT ["/usr/local/bin/zenx"]
