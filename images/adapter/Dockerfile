FROM golang:1.13 as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY config/ ./config
COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux \
    go build -a -installsuffix cgo -v -o bin/kuda_adapter .

# Second stage

FROM alpine:3.8
RUN apk --no-cache add ca-certificates

WORKDIR /bin/

COPY --from=builder /app/bin/kuda_adapter .

ENTRYPOINT [ "/bin/kuda_adapter" ]

CMD [ "44225" ]
EXPOSE 44225