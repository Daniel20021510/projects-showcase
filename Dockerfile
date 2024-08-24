FROM golang:1.23-alpine AS builder

WORKDIR /usr/local/src/projectsShowcase

RUN apk --no-cache add gcc musl-dev

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal
RUN go build -o projectsShowcase cmd/projectsShowcase/main.go

FROM alpine AS runner

COPY --from=builder /usr/local/src/projectsShowcase/projectsShowcase /
COPY /config/local.yaml /config.yaml
RUN mkdir "storage"

CMD ["env","CONFIG_PATH=/config.yaml","/projectsShowcase"]