FROM golang:1.18.2-alpine
ENV POSTGRESURI="postgresql://test:test@db:5432/sharesDB?sslmode=disable"
ENV IMPORTCSV="YES"
WORKDIR /opt
RUN apk add gcc libc-dev libpcap-dev
COPY . .
#RUN go mod tidy
RUN go build -mod vendor ./cmd/dhtcrwl/...
EXPOSE 8080
EXPOSE 1337
USER root
CMD ["./dhtcrwl"]
