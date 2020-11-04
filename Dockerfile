FROM golang:alpine AS builder

WORKDIR /GO-BD-MYSQL

COPY . .

RUN  go mod download

RUN go mod verify


RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/exec

FROM scratch

COPY --from=builder /go/bin/exec /go/bin/exec

EXPOSE 8080

CMD ["/go/bin/exec"]

