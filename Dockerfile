FROM golang:1.20 AS builder
COPY . /src 
WORKDIR /src
ENV CGO_ENABLED=0 
RUN go build -o app ./cmd/scoreboard_exporter

FROM scratch AS final 
LABEL org.opencontainers.image.source=https://github.com/boxmein/adctf_scoreboard_exporter
COPY --from=builder /src/app /app 
ENTRYPOINT ["/app"]
