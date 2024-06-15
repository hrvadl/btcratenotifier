FROM golang:latest AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod tidy
COPY ./ ./
RUN cd cmd/load && CGO_ENABLED=0 GOOS=linux go build .

FROM scratch AS final
WORKDIR /app
COPY --from=builder /app/cmd/load/load .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
CMD [ "./load" ]
