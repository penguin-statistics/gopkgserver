FROM golang:1.27.3-alpine AS base
WORKDIR /app

# builder
FROM base AS builder
ENV GOOS linux
ENV GOARCH amd64

# build-args
ARG VERSION

RUN apk --no-cache add bash git openssh

# modules: utilize build cache
COPY go.mod ./
COPY go.sum ./

# RUN go env -w GO111MODULE=on && go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod download
COPY . .

# build the binary (static, so the runtime image needs no Go toolchain or libc shims)
RUN CGO_ENABLED=0 go build -trimpath -o gopkgserver .

# runner: plain Alpine, no Go toolchain
FROM alpine:3.24 AS runner
WORKDIR /app
RUN apk add --no-cache tini
# Tini is now available at /sbin/tini

COPY --from=builder /app/gopkgserver /app/gopkgserver
EXPOSE 3000

ENTRYPOINT ["/sbin/tini", "--"]
CMD [ "/app/gopkgserver" ]
