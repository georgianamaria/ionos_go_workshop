FROM golang:1.21-alpine AS builder

WORKDIR /src

COPY . /src

RUN go build -o /bin/service

FROM gcr.io/distroless/base-debian11:nonroot

COPY --from=builder /bin/service /bin/service

EXPOSE 8080

ENTRYPOINT ["/bin/service"]