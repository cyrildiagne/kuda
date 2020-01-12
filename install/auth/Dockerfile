FROM golang:1.13 as builder

# Create and change to the app directory.
WORKDIR /app

# Retrieve application dependencies using go modules.
# Allows container builds to reuse downloaded dependencies.
COPY go.* ./
RUN go mod download

COPY pkg ./pkg
COPY cmd/auth ./cmd/auth
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -mod=readonly -installsuffix cgo -o auth ./cmd/auth

COPY web/auth ./web/auth

#

FROM alpine:3
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/auth ./auth

# Copy public assets to the container image.
COPY  --from=builder /app/web/auth ./web/auth

# Run the web service on container startup.
CMD ["/auth"]