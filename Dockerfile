FROM golang:1.13
WORKDIR /go/src/github.com/thomasfady/httpfiletransfer/
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./httpfiletransfer

FROM alpine:latest  
WORKDIR /root/
COPY --from=0 /go/src/github.com/thomasfady/httpfiletransfer/httpfiletransfer .
RUN mkdir uploads static tpl
ENTRYPOINT ["./httpfiletransfer"]  