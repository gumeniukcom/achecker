FROM golang:1.15.0-alpine@sha256:73182a0a24a1534e31ad9cc9e3a4bb46bb030a883b26eda0a87060f679b83607 as builder
RUN apk add --update alpine-sdk
###

WORKDIR /go/src/githiub.com/gumeniukcom/achekcer
COPY . .
RUN CGO_ENABLED=0 GOOS=linux make all

###
FROM alpine:3.12.0@sha256:185518070891758909c9f839cf4ca393ee977ac378609f700f60a771a2dfe321
RUN apk --no-cache add tzdata

RUN adduser -D achecker
USER achecker
WORKDIR /home/achecker/

ENV TZ=UTC

COPY --from=builder --chown=achecker:achecker /go/src/githiub.com/gumeniukcom/achekcer/achecker .
COPY --from=builder --chown=achecker:achecker /go/src/githiub.com/gumeniukcom/achekcer/config.toml .
RUN mkdir -p  /home/achecker/kafkaca
COPY --chown=achecker:achecker kafkaca kafkaca

###
LABEL maintainer="Stanislav Gumeniuk i@gumeniuk.com"
###

CMD ./achecker
