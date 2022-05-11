FROM golang
ENV POSTGRESURI="postgresql://test:test@db:5432/sharesDB?sslmode=disable"
ENV THREADS="5"
WORKDIR /opt
COPY . .
#RUN go mod tidy
RUN go build -mod vendor ./cmd/dhtcrwl/...
EXPOSE 8080
CMD ["./dhtcrwl"]
