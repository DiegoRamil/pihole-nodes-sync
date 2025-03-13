FROM -platform=${BUILDPLATFORM:-linux/amd64} golang:1.23-alpine AS build

ARG GOARCH
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH


WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${GOARCH} go build -ldflags="-w -s" -o main cmd/app/main.go

FROM --platform=${TARGETPLATFORM:-linux/amd64} scratch AS prod
WORKDIR /app
COPY --from=build /app/main /app/main
EXPOSE ${PORT}
CMD ["./main"]


