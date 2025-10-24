FROM golang:1.25-alpine AS BUILDER

ENV CGO_ENABLED=0
WORKDIR /myapp

COPY go.mod go.sum ./
RUN go mod download

RUN make build

FROM scratch
COPY --from=BUILDER /myapp/myapp /myapp
CMD ["/myapp"]
