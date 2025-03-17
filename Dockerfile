FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.23.4 AS builder

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o main cmd/app/main.go

FROM --platform=${BUILDPLATFORM:-linux/amd64} scratch AS prod
WORKDIR /app
COPY --from=builder /app/main /app/main
EXPOSE ${PORT}
CMD ["./main"]
