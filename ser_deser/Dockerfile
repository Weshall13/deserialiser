FROM golang:latest as builder
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /ser_deser

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /ser_deser .
COPY templates/ ./templates/
COPY static/ ./static/
EXPOSE 8080
CMD ["/root/ser_deser"] 