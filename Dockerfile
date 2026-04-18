FROM alpine:3.14

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root

# Copy binary đã build sẵn và assets
COPY permission .

EXPOSE 7000 7001

CMD ["./permission", "start"]