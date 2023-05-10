FROM golang:1.20-alpine AS build
WORKDIR /home
COPY  . .
RUN go build -o exec

FROM alpine
WORKDIR /home
COPY --from=build /home/exec /home/app

ARG DEFAULT_NATS_URL "localhost:4222"
ENV NATS_URL $DEFAULT_NATS_URL

CMD ["./app"]