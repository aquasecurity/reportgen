FROM golang:1.17.3-alpine as builder
RUN mkdir /reportgen
ADD . /reportgen/
WORKDIR /reportgen
RUN go build -o main .

FROM alpine
RUN mkdir /reportgen
RUN mkdir /reportgen/assets
RUN mkdir /reports
COPY --from=builder /reportgen/main /reportgen/
COPY --from=builder /reportgen/assets/*.ttf /reportgen/assets/
COPY --from=builder /reportgen/assets/*.png /reportgen/assets/
WORKDIR /reportgen
RUN addgroup -g 1099 report
RUN adduser -D -g '' -G report -u 1099 report
RUN chown -R report:report /reportgen
RUN chown -R report:report /reports
USER report
VOLUME ["/reports"]
ENTRYPOINT ["/reportgen/main"]
