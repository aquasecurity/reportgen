FROM golang:alpine as builder
RUN apk add git
RUN go get github.com/signintech/gopdf
RUN mkdir /reportgen
ADD . /reportgen/
WORKDIR /reportgen
RUN go build -o main .

FROM alpine
RUN mkdir /reportgen
RUN mkdir /reportgen/pdfrender
COPY --from=builder /reportgen/main /reportgen/
COPY --from=builder /reportgen/pdfrender/*.ttf /reportgen/pdfrender/
COPY --from=builder /reportgen/pdfrender/*.png /reportgen/pdfrender/
WORKDIR /reportgen
#RUN adduser -S -D -H -h /reportgen report
#USER report
ENTRYPOINT ["/reportgen/main"]
