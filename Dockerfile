FROM docker.io/golang:1.17-alpine as builder

RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build .

FROM docker.io/alpine:latest
COPY --from=builder /app/ot2mqtt /
CMD ["/ot2mqtt"]
