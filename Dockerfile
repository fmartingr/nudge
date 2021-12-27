FROM golang:1.16.12-alpine3.15 AS builder

RUN apk update && \
    apk add --no-cache --update make bash git ca-certificates && \
    update-ca-certificates

WORKDIR /app
ENV GO111MODULE=on
COPY . .

RUN rm -rf build && mkdir build
RUN make build

FROM alpine:3.13

COPY docker/entry.sh ./
COPY --from=builder /app/build/* /usr/local/bin/

ENTRYPOINT [ "./entry.sh" ]

CMD [ "nudge" ]
