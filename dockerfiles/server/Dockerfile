FROM golang:1.19.4 AS builder
WORKDIR /
COPY ./server ./
RUN rm -rf vendor

RUN go mod download
RUN CGO_ENABLED=0 go build -o ./build/server cmd/server/main.go

FROM scratch
COPY --from=builder ./build/server server

USER 65532:65532
ENTRYPOINT ["./server"]
