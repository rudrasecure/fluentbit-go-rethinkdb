FROM golang

ENV GOOS=linux\
    GOARCH=amd64

WORKDIR /
RUN git clone https://github.com/rudrasecure/fluentbit-go-rethinkdb
WORKDIR /fluentbit-go-rethinkdb
RUN go build -o out/fluentbit-go-rethinkdb.so -buildmode=c-shared

FROM fluent/fluent-bit

COPY --from=0 /fluentbit-go-rethinkdb/out/fluentbit-go-rethinkdb.so /fluent-bit/bin/out/fluentbit-go-rethinkdb.so
COPY fluent-bit.conf /fluent-bit/etc/fluent-bit.conf
COPY plugins.conf /fluent-bit/etc/plugins.conf

CMD ["/fluent-bit/bin/fluent-bit", "-c", "/fluent-bit/etc/fluent-bit.conf"]