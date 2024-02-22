# Get base Go image
FROM golang:1.22-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o frontendApp ./cmd/web

RUN chmod +x  /app/frontendApp

# # build a tiny docker image

FROM alpine:latest

RUN mkdir /app

RUN mkdir /templates

COPY --from=builder /app/frontendApp /app/frontendApp

COPY ./cmd/web/templates/ /templates

CMD ["/app/frontendApp"]
