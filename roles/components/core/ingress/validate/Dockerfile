# must match spec.volumeMounts[0].mountPath in sonobuoy config
ARG resultsDir="/tmp/results"

# build
FROM golang:1.15 AS builder
WORKDIR /test
COPY . /test/
RUN CGO_ENABLED=0 go build -o /bin/run ./main.go

# runtime
FROM alpine:3.12
WORKDIR /test

RUN addgroup -S nonroot && adduser -S nonroot -G nonroot

COPY --from=builder --chown=nonroot:nonroot /bin/run /test/run

USER nonroot:nonroot
CMD ["/test/run"]

