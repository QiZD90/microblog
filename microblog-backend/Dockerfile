# Cache modules
FROM golang:alpine as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Build the app
FROM golang:alpine AS build
WORKDIR /build
COPY --from=modules /go/pkg /go/pkg
COPY . .
RUN go build -o microblog-backend cmd/app/main.go

# Run the app
FROM alpine
WORKDIR /app
COPY --from=build /build/microblog-backend .
CMD ["./microblog-backend"]