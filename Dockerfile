FROM golang:1.10 as go

COPY . /go/src/law
WORKDIR /go/src/law

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o law *.go

FROM alpine:3.6
RUN apk --no-cache add ca-certificates
COPY --from=go /go/src/law/law /app/law
COPY files/server.crt /app/server.crt
COPY files/server.key /app/server.key

EXPOSE 80
CMD ["/app/law"]
