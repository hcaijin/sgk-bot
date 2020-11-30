FROM golang:1.14 as builder

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...

RUN make clean \
    && make linux

FROM alpine

WORKDIR /app

COPY --from=builder /go/src/app/bin/sgkbot_linux .

RUN set -ex \
      && apk update \
      && apk add --update tzdata

ENV TZ=Asia/Shanghai

RUN rm -rf /var/cache/apk/*

CMD [ "./sgkbot_linux" ]
