FROM qdtech.tencentcloudcr.com/base/go-builder:0.0.1 as builder

WORKDIR /www

COPY . .

RUN make build

FROM qdtech.tencentcloudcr.com/base/go-runner:0.0.1 as runner

ENV GIN_MODE release

WORKDIR /www

COPY --from=builder /www/bin/runner /www/eztalk-monitor

CMD [ "/www/eztalk-monitor" ]

