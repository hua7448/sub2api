FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY image/sub2api /app/sub2api
RUN chmod 0755 /app/sub2api && mkdir -p /app/data
EXPOSE 18082
ENTRYPOINT ["/app/sub2api"]
