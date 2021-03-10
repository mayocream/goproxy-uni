FROM golang:1.16

WORKDIR /gu

COPY go.mod .

COPY go.sum .

RUN go env -w GOPROXY=https://proxy.golang.org,direct \
    && go mod download

COPY . .

RUN go build -o goproxy_uni

RUN touch .netrc \
    && ln -sf $(pwd)/.netrc ~/.netrc

CMD ["./goproxy_uni", "-c", "config.yaml"]