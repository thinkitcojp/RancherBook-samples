FROM golang:1.11 as builder
RUN go get -u github.com/spf13/cobra/cobra
COPY gitlab.com /go/src/gitlab.com
WORKDIR /go/src/gitlab.com/fufuhu/ti_rancher_k8s_sampleapp
# cmdパッケージのテストバイナリのビルドとテスト
# RUN cd ./cmd && go test ./ -c  && cd ..
# serviceパッケージのテストバイナリのビルドとテスト
# RUN cd ./service && go test ./ -c && cd ..
# クライアントアプリケーション本体のビルド
RUN go build -o todo

FROM alpine:3.9
COPY --from=builder /go/src/gitlab.com/fufuhu/ti_rancher_k8s_sampleapp/todo ./