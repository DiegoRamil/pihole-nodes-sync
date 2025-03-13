FROM golang:1.23-alpine AS build
ARG GO_ARCH

WORKDIR /app

COPY go.mod go.sum ./
RUN GOARCH=${GO_ARCH} go mod download

COPY . .

RUN GOARCH=${GO_ARCH} go build -o main cmd/app/main.go

FROM scratch AS prod
WORKDIR /app
COPY --from=build /app/main /app/main
EXPOSE ${PORT}
CMD ["./main"]


