FROM golang:1.22 AS base

WORKDIR /app

COPY go-web-app/go.mod ./

RUN go mod download

COPY go-web-app/. .

RUN go build -o main .

# minimal distroless image

FROM gcr.io/distroless/base

COPY --from=base /app/main .

COPY --from=base /app/static ./static

EXPOSE 8080

USER 65532:65532

CMD ["./main"]