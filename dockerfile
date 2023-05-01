FROM golang:1.18-alpine AS build

WORKDIR /app
COPY go.mod .
COPY . .
RUN go build -o main ./cmd

# copy binary into smaller base image
FROM alpine:3.14
WORKDIR /app
COPY --from=build /app .

# docker compose wait tool
ENV WAIT_VERSION 2.7.2
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/$WAIT_VERSION/wait /wait
RUN chmod +x /wait

EXPOSE 8080

CMD ["./main"]