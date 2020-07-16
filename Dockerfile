# docker build -t luanpeng/lp:image-go-1.0.0 .

FROM golang

RUN go get github.com/astaxie/beego
RUN go get github.com/prometheus/client_golang/prometheus/promhttp
RUN go get github.com/prometheus/client_golang/prometheus
COPY conf /go/src/image_go/conf
COPY controllers /go/src/image_go/controllers
COPY main.go /go/src/image_go/main.go

WORKDIR /go/src/image_go/
RUN go build main.go
CMD ["./main"]








