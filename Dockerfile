FROM golang:1.14.0-stretch as builder
RUN mkdir /src
WORKDIR /src
COPY * /src
RUN CGO_ENABLED=0 go build -o /app .

FROM gcr.io/distroless/static:48dba0a4ace4fcb4fdd8d7e1f7dc1a9ed8b38f7c
COPY --from=builder /app /
CMD ["/app"]
