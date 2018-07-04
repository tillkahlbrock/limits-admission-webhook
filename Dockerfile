FROM golang:1.10 as go

COPY . /go/src/law
WORKDIR /go/src/law

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o law *.go

FROM alpine:3.6
RUN apk --no-cache add ca-certificates
COPY --from=go /go/src/law/law /bin/law

EXPOSE 80
CMD ["/bin/law"]
