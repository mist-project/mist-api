# -------- Build stage --------
FROM golang:1.23-alpine AS builder

ARG APP_PORT=4000

ENV CGO_ENABLED=0 \
    GO111MODULE=on \
    PROJECT_DIR=/app

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Install swag CLI for Swagger docs generation
RUN go install github.com/swaggo/swag/cmd/swag@v1.16.6

# Copy only Go source files first (exclude docs and bin for caching)
COPY ./src ./src

# Generate Swagger docs
RUN swag init --dir src

# Build binary into bin/ directory
RUN mkdir -p bin && go build -o bin/mistapi src/main.go

# -------- Final stage --------
FROM gcr.io/distroless/static:nonroot

ARG APP_PORT
ENV APP_PORT=${APP_PORT}
EXPOSE ${APP_PORT}

WORKDIR /app

COPY --from=builder /app/bin/mistapi .

USER nonroot:nonroot

ENTRYPOINT ["/app/mistapi"]
