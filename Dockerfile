FROM golang:1.26-alpine AS build
WORKDIR /app
COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal
RUN go build -o /shipwatch ./cmd/shipwatch

FROM alpine:3.22
COPY --from=build /shipwatch /usr/local/bin/shipwatch
EXPOSE 8080
CMD ["shipwatch"]
