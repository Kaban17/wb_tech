FROM golang:1.25-alpine AS builder

ENV CGO_ENABLED=0
WORKDIR /myapp

COPY . .
RUN go mod download

RUN go build -o myapp ./cmd/app

FROM scratch
# Set the working directory for the final image
WORKDIR /app

# Copy the compiled binary to the new working directory
COPY --from=builder /myapp/myapp /app/myapp

# Copy the .env file to the working directory
COPY .env .env

# Copy the web directory to the root of the container, as it's referenced by http.Dir("/web")
# This path is absolute and not dependent on WORKDIR.
COPY web/ /web/

# Set the command to run the application
CMD ["/app/myapp"]