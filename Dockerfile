FROM golang:alpine as builder
RUN apk add git
RUN go get github.com/signintech/gopdf
RUN go get github.com/joho/godotenv
RUN mkdir /reportgen
ADD . /reportgen/
WORKDIR /reportgen
RUN go build -o main .

FROM alpine
RUN mkdir /reportgen
RUN mkdir /reportgen/pdfrender
RUN mkdir /reports
COPY --from=builder /reportgen/main /reportgen/
COPY --from=builder /reportgen/pdfrender/*.ttf /reportgen/pdfrender/
COPY --from=builder /reportgen/pdfrender/*.png /reportgen/pdfrender/
WORKDIR /reportgen
RUN adduser -D -g '' report
RUN chown -R report:report /reportgen
RUN chown -R report:report /reports
USER report
VOLUME ["/reports"]
ENTRYPOINT ["/reportgen/main"]
