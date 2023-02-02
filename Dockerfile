FROM golang:1.19.1

WORKDIR /opt/fazz

COPY ./app /opt/fazz

RUN go build -o /app/build/app ./cmd/app/main.go
CMD /app/build/app
