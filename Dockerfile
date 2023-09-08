# Build phase
FROM golang:alpine AS build

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o service .

# Run phase
FROM alpine:latest
ARG RAFT_VOL_DIR

RUN apk add --update-cache \
    iproute2 \
    iputils-ping

WORKDIR /app
COPY --from=build /app/service .

RUN mkdir -p $RAFT_VOL_DIR
CMD [ "./service" ]
